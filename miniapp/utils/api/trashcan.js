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

function createTrashCanWithImages(formData, imagePaths) {
  return new Promise((resolve, reject) => {
    const token = wx.getStorageSync('token')
    if (!token) {
      reject({ message: '未登录' })
      return
    }

    const uploadUrl = `${API_ORIGIN}/api/trashcans`

    let uploadedPaths = []
    let currentIndex = 0

    const uploadNext = () => {
      if (currentIndex >= imagePaths.length) {
        const data = { ...formData }
        if (uploadedPaths[0]) data.image = uploadedPaths[0]
        if (uploadedPaths[1]) data.image_2 = uploadedPaths[1]
        if (uploadedPaths[2]) data.image_3 = uploadedPaths[2]

        wx.request({
          url: uploadUrl,
          method: 'POST',
          header: {
            'Authorization': `Bearer ${token}`,
            'content-type': 'application/x-www-form-urlencoded'
          },
          data,
          success: (res) => {
            if (res.data && res.data.code === 0) {
              resolve(res.data)
            } else {
              reject({ message: res.data?.msg || '创建失败' })
            }
          },
          fail: (err) => {
            reject({ message: err.errMsg || '请求失败' })
          }
        })
        return
      }

      wx.uploadFile({
        url: uploadUrl,
        filePath: imagePaths[currentIndex],
        name: `images[${currentIndex}]`,
        header: {
          'Authorization': `Bearer ${token}`
        },
        formData: currentIndex === 0 ? formData : {},
        success: (res) => {
          if (res.statusCode === 200) {
            try {
              const data = JSON.parse(res.data)
              if (data.code === 0 && data.data && data.data.image_path) {
                uploadedPaths.push(data.data.image_path)
              }
            } catch (e) {}
          }
          currentIndex++
          uploadNext()
        },
        fail: (err) => {
          currentIndex++
          uploadNext()
        }
      })
    }

    uploadNext()
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
  createTrashCanWithImages,
  getUserTrashCans,
  toggleLike,
  toggleDislike,
  deleteTrashCan,
  updateTrashCan
}
