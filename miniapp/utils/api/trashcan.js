const { request } = require("../request")
const { API_ORIGIN } = require("../../config/index")

function getNearbyTrashCans(lat, lng, radius = 5, limit = 10, tagId = '') {
  return request({
    url: "/trashcans/nearby",
    method: "GET",
    data: { lat, lng, radius, limit, tag_id: tagId }
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
      wx.showToast({ title: '未登录', icon: 'none' })
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
            if ((json.code === 2000 || json.code === 0) && json.data && json.data.image_path) {
              resolve(json.data.image_path)
            } else {
              wx.showToast({ title: json.msg || '上传图片失败', icon: 'none' })
              reject({ message: json.msg || '上传图片失败' })
            }
          } catch (e) {
            wx.showToast({ title: '解析响应失败', icon: 'none' })
            reject({ message: '解析响应失败' })
          }
        } else {
          wx.showToast({ title: '上传失败: ' + res.statusCode, icon: 'none' })
          reject({ message: '上传失败' })
        }
      },
      fail: (err) => {
        wx.showToast({ title: '上传失败: ' + err.errMsg, icon: 'none' })
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

    const uploadPromises = imagePaths.map(path => 
      uploadImage(path).catch(err => {
        return ''
      })
    )
    
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

function getAdminTrashCans(page = 1, pageSize = 20, keyword = '') {
  return request({
    url: "/admin/trashcans",
    method: "GET",
    data: { page, page_size: pageSize, keyword }
  })
}

function updateAdminTrashCan(id, data) {
  return request({
    url: `/admin/trashcans/${id}`,
    method: "PUT",
    data: data,
    header: {
      "content-type": "application/x-www-form-urlencoded"
    }
  })
}

function deleteAdminTrashCan(id) {
  return request({
    url: `/admin/trashcans/${id}`,
    method: "DELETE"
  })
}

function getAdminTags() {
  return request({
    url: "/admin/tags",
    method: "GET"
  })
}

function deleteAdminTag(id) {
  return request({
    url: "/admin/tags",
    method: "DELETE",
    data: { id }
  })
}

function createAdminTag(name) {
  return request({
    url: "/admin/tags",
    method: "POST",
    data: { name },
    header: {
      "content-type": "application/x-www-form-urlencoded"
    }
  })
}

function getAllTags() {
  return request({
    url: "/tags",
    method: "GET"
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
  updateTrashCan,
  getAdminTrashCans,
  updateAdminTrashCan,
  deleteAdminTrashCan,
  getAdminTags,
  deleteAdminTag,
  createAdminTag,
  getAllTags
}
