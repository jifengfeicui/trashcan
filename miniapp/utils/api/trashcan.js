const { request } = require("../request")

function getNearbyTrashCans(lat, lng, radius = 5, limit = 10) {
  return request({
    url: "/trashcans/nearby",
    method: "GET",
    data: { lat, lng, radius, limit }
  })
}

function createTrashCan(formData) {
  return request({
    url: "/trashcans",
    method: "POST",
    data: formData,
    header: {
      "content-type": "application/x-www-form-urlencoded"
    }
  })
}

function createTrashCanWithImages(formData, imagePaths) {
  return new Promise((resolve, reject) => {
    const token = wx.getStorageSync('token')
    if (!token) {
      reject({ message: '未登录' })
      return
    }

    const uploadDir = `${wx.env.HTTP_BASE}/api/trashcans`
    const header = {
      'Authorization': `Bearer ${token}`
    }

    let uploadedCount = 0
    const uploadedPaths = []

    if (imagePaths.length === 0) {
      wx.request({
        url: uploadDir,
        method: 'POST',
        header,
        data: formData,
        success: (res) => {
          if (res.data.code === 0) {
            resolve(res.data)
          } else {
            reject({ message: res.data.msg || '上传失败' })
          }
        },
        fail: reject
      })
      return
    }

    const uploadNext = (index) => {
      if (index >= imagePaths.length) {
        formData.images = uploadedPaths
        wx.request({
          url: uploadDir,
          method: 'POST',
          header,
          data: formData,
          success: (res) => {
            if (res.data.code === 0) {
              resolve(res.data)
            } else {
              reject({ message: res.data.msg || '上传失败' })
            }
          },
          fail: reject
        })
        return
      }

      wx.uploadFile({
        url: uploadDir,
        filePath: imagePaths[index],
        name: `images[${index}]`,
        header,
        formData: index === 0 ? formData : {},
        success: (res) => {
          try {
            const data = JSON.parse(res.data)
            if (data.code === 0) {
              uploadedPaths.push(data.data.image_path || '')
              uploadedCount++
              uploadNext(index + 1)
            } else {
              reject({ message: data.msg || '上传失败' })
            }
          } catch (e) {
            if (res.statusCode === 200) {
              uploadedPaths.push('')
              uploadedCount++
              uploadNext(index + 1)
            } else {
              reject({ message: '上传失败' })
            }
          }
        },
        fail: (err) => {
          reject(err)
        }
      })
    }

    uploadNext(0)
  })
}

function getUserTrashCans(page = 1, pageSize = 10) {
  return request({
    url: "/users/me/trashcans",
    method: "GET",
    data: {
      page,
      page_size: pageSize
    }
  })
}

function toggleLike(id) {
  return request({
    url: `/trashcans/${id}/like`,
    method: "POST"
  })
}

function toggleDislike(id) {
  return request({
    url: `/trashcans/${id}/dislike`,
    method: "POST"
  })
}

function deleteTrashCan(id) {
  return request({
    url: `/trashcans/${id}`,
    method: "DELETE"
  })
}

module.exports = {
  getNearbyTrashCans,
  createTrashCan,
  createTrashCanWithImages,
  getUserTrashCans,
  toggleLike,
  toggleDislike,
  deleteTrashCan
}
