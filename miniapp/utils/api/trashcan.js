const { request } = require("../request")
const { API_ORIGIN } = require("../../config/index")

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

function uploadImage(filePath) {
  return new Promise((resolve, reject) => {
    const token = wx.getStorageSync('token')
    if (!token) {
      reject({ message: '未登录' })
      return
    }

    const uploadUrl = `${API_ORIGIN}/api/upload/image`
    
    wx.uploadFile({
      url: uploadUrl,
      filePath: filePath,
      name: 'image',
      header: {
        'Authorization': `Bearer ${token}`
      },
      success: (res) => {
        if (res.statusCode === 200) {
          try {
            const json = JSON.parse(res.data)
            if (json.code === 0 && json.data && json.data.image_path) {
              resolve(json.data.image_path)
            } else {
              reject({ message: json.msg || '上传图片失败' })
            }
          } catch (e) {
            reject({ message: '解析响应失败' })
          }
        } else {
          reject({ message: '上传失败' })
        }
      },
      fail: (err) => {
        reject({ message: err.errMsg || '上传失败' })
      }
    })
  })
}

function createTrashCanWithImages(formData, imagePaths) {
  return new Promise((resolve, reject) => {
    const token = wx.getStorageSync('token')
    if (!token) {
      reject({ message: '未登录' })
      return
    }

    if (!imagePaths || imagePaths.length === 0) {
      createTrashCan(formData).then(resolve).catch(reject)
      return
    }

    const uploadPromises = imagePaths.map(path => uploadImage(path))
    
    Promise.all(uploadPromises)
      .then(imagePathsResult => {
        const data = { ...formData }
        if (imagePathsResult[0]) data.image = imagePathsResult[0]
        if (imagePathsResult[1]) data.image_2 = imagePathsResult[1]
        if (imagePathsResult[2]) data.image_3 = imagePathsResult[2]

        createTrashCan(data).then(resolve).catch(reject)
      })
      .catch(err => {
        reject(err)
      })
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

function updateTrashCan(id, data) {
  return request({
    url: `/trashcans/${id}`,
    method: "PUT",
    data: data,
    header: {
      "content-type": "application/x-www-form-urlencoded"
    }
  })
}

module.exports = {
  getNearbyTrashCans,
  createTrashCan,
  uploadImage,
  createTrashCanWithImages,
  getUserTrashCans,
  toggleLike,
  toggleDislike,
  deleteTrashCan,
  updateTrashCan
}
