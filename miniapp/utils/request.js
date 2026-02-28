const { API_BASE_URL, API_ORIGIN } = require("../config/index")
const { getToken, clearAuth } = require("./auth")

function handleResponseBody(body, resolve, reject) {
  if (!body || body.code === undefined) {
    resolve(body)
    return
  }

  if (body.code === 2000) {
    resolve(body)
    return
  }

  reject(new Error(body.msg || "请求失败"))
}

function request(options) {
  const {
    url,
    method = "GET",
    data = {},
    header = {},
    timeout = 10000
  } = options

  return new Promise((resolve, reject) => {
    const token = getToken()
    const mergedHeader = Object.assign({}, header)

    if (token) {
      mergedHeader.Authorization = `Bearer ${token}`
    }

    wx.request({
      url: `${API_BASE_URL}${url}`,
      method,
      data,
      header: mergedHeader,
      timeout,
      success: (res) => {
        if (res.statusCode === 401) {
          clearAuth()
          reject(new Error("未登录或登录已过期"))
          return
        }

        handleResponseBody(res.data, resolve, reject)
      },
      fail: (err) => {
        reject(new Error(err.errMsg || "网络请求失败"))
      }
    })
  })
}

function uploadFile(options) {
  const {
    url,
    filePath,
    name = "image",
    formData = {},
    header = {}
  } = options

  return new Promise((resolve, reject) => {
    const token = getToken()
    const mergedHeader = Object.assign({}, header)

    if (token) {
      mergedHeader.Authorization = `Bearer ${token}`
    }

    wx.uploadFile({
      url: `${API_BASE_URL}${url}`,
      filePath,
      name,
      formData,
      header: mergedHeader,
      success: (res) => {
        if (res.statusCode === 401) {
          clearAuth()
          reject(new Error("未登录或登录已过期"))
          return
        }

        let body = null
        try {
          body = JSON.parse(res.data)
        } catch (e) {
          reject(new Error("服务端返回格式错误"))
          return
        }

        handleResponseBody(body, resolve, reject)
      },
      fail: (err) => {
        reject(new Error(err.errMsg || "上传失败"))
      }
    })
  })
}

function resolveImageURL(imageURL) {
  if (!imageURL) {
    return ""
  }

  if (/^https?:\/\//i.test(imageURL)) {
    return imageURL
  }

  if (imageURL.startsWith("/")) {
    return `${API_ORIGIN}${imageURL}`
  }

  return `${API_ORIGIN}/${imageURL}`
}

module.exports = {
  request,
  uploadFile,
  resolveImageURL
}
