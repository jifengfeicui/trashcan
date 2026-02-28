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
  getUserTrashCans,
  toggleLike,
  toggleDislike,
  deleteTrashCan
}
