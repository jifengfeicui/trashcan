const { AMAP_WEB_KEY } = require("../config/index")

function getAddressByLocation(lng, lat) {
  if (!lng || !lat) {
    return Promise.resolve("")
  }

  return new Promise((resolve, reject) => {
    wx.request({
      url: "https://restapi.amap.com/v3/geocode/regeo",
      method: "GET",
      data: {
        key: AMAP_WEB_KEY,
        location: `${lng},${lat}`,
        radius: 1000,
        extensions: "all"
      },
      success: (res) => {
        const body = res.data || {}
        if (body.status === "1" && body.regeocode) {
          resolve(body.regeocode.formatted_address || "")
          return
        }

        reject(new Error(body.info || "逆地理编码失败"))
      },
      fail: (err) => {
        reject(new Error(err.errMsg || "逆地理编码失败"))
      }
    })
  })
}

module.exports = {
  getAddressByLocation
}
