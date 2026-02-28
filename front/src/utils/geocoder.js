/**
 * 高德地图逆地理编码工具
 * 根据经纬度获取地址描述
 * 使用高德地图Web API
 */

const AMAP_WEB_KEY = '811010fd93771b8d6b94c534eebe4133'
const AMAP_REGEOCODE_API = 'https://restapi.amap.com/v3/geocode/regeo'

/**
 * 调用高德地图逆地理编码API
 * @param {number} lng - 经度
 * @param {number} lat - 纬度
 * @returns {Promise<Object>} 返回API响应数据
 */
const callRegeocodeAPI = async (lng, lat) => {
  const url = `${AMAP_REGEOCODE_API}?key=${AMAP_WEB_KEY}&location=${lng},${lat}&radius=1000&extensions=all`
  
  try {
    const response = await fetch(url)
    const data = await response.json()
    
    if (data.status === '1' && data.regeocode) {
      return data
    } else {
      throw new Error(data.info || '逆地理编码失败')
    }
  } catch (error) {
    throw new Error(`逆地理编码请求失败: ${error.message}`)
  }
}

/**
 * 根据经纬度获取地址描述
 * @param {number} lng - 经度
 * @param {number} lat - 纬度
 * @returns {Promise<string>} 返回地址描述字符串
 */
export const getAddressByLocation = async (lng, lat) => {
  try {
    const data = await callRegeocodeAPI(lng, lat)
    return data.regeocode.formatted_address || ''
  } catch (error) {
    console.error('逆地理编码错误:', error)
    throw error
  }
}

/**
 * 根据经纬度获取详细地址信息
 * @param {number} lng - 经度
 * @param {number} lat - 纬度
 * @returns {Promise<Object>} 返回详细地址信息对象
 */
export const getAddressDetailByLocation = async (lng, lat) => {
  try {
    const data = await callRegeocodeAPI(lng, lat)
    const regeocode = data.regeocode
    const addressComponent = regeocode.addressComponent || {}
    
    return {
      formattedAddress: regeocode.formatted_address || '',
      province: addressComponent.province || '',
      city: addressComponent.city || '',
      district: addressComponent.district || '',
      township: addressComponent.township || '',
      street: addressComponent.street || '',
      streetNumber: addressComponent.streetNumber || '',
      // 只使用 formattedAddress，不再手动拼接
      address: regeocode.formatted_address || ''
    }
  } catch (error) {
    console.error('逆地理编码错误:', error)
    throw error
  }
}

