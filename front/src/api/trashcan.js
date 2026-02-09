import request from '@/utils/request'

/**
 * 获取附近的垃圾桶
 * @param {number} lat - 纬度
 * @param {number} lng - 经度
 * @param {number} radius - 搜索半径（公里），默认5
 * @param {number} limit - 返回数量限制，默认10
 * @returns {Promise}
 */
export function getNearbyTrashCans(lat, lng, radius = 5, limit = 10) {
  return request({
    url: '/trashcans/nearby',
    method: 'get',
    params: {
      lat,
      lng,
      radius,
      limit
    }
  })
}

/**
 * 创建垃圾桶
 * @param {FormData} formData - 表单数据，包含latitude, longitude, address, description, image
 * @returns {Promise}
 */
export function createTrashCan(formData) {
  return request({
    url: '/trashcans',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

/**
 * 获取垃圾桶详情
 * @param {number|string} id - 垃圾桶ID
 * @returns {Promise}
 */
export function getTrashCanDetail(id) {
  return request({
    url: `/trashcans/${id}`,
    method: 'get'
  })
}

/**
 * 点赞垃圾桶
 * @param {number|string} id - 垃圾桶ID
 * @returns {Promise}
 */
export function toggleLike(id) {
  return request({
    url: `/trashcans/${id}/like`,
    method: 'post'
  })
}

/**
 * 点踩垃圾桶
 * @param {number|string} id - 垃圾桶ID
 * @returns {Promise}
 */
export function toggleDislike(id) {
  return request({
    url: `/trashcans/${id}/dislike`,
    method: 'post'
  })
}

/**
 * 获取用户上传的垃圾桶列表（分页）
 * @param {number} page - 页码，默认1
 * @param {number} pageSize - 每页数量，默认10
 * @returns {Promise}
 */
export function getUserTrashCans(page = 1, pageSize = 10) {
  return request({
    url: '/users/me/trashcans',
    method: 'get',
    params: {
      page,
      page_size: pageSize
    }
  })
}

/**
 * 更新垃圾桶信息
 * @param {number|string} id - 垃圾桶ID
 * @param {FormData} formData - 表单数据，包含address, description, image（可选）
 * @returns {Promise}
 */
export function updateTrashCan(id, formData) {
  return request({
    url: `/trashcans/${id}`,
    method: 'put',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

/**
 * 删除垃圾桶
 * @param {number|string} id - 垃圾桶ID
 * @returns {Promise}
 */
export function deleteTrashCan(id) {
  return request({
    url: `/trashcans/${id}`,
    method: 'delete'
  })
}

