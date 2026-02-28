const TOKEN_KEY = "token"
const USER_KEY = "user_info"

function getToken() {
  return wx.getStorageSync(TOKEN_KEY) || ""
}

function setToken(token) {
  wx.setStorageSync(TOKEN_KEY, token || "")
}

function getUserInfo() {
  return wx.getStorageSync(USER_KEY) || null
}

function setUserInfo(userInfo) {
  wx.setStorageSync(USER_KEY, userInfo || null)
}

function clearAuth() {
  wx.removeStorageSync(TOKEN_KEY)
  wx.removeStorageSync(USER_KEY)
}

function isAuthenticated() {
  return !!getToken()
}

module.exports = {
  getToken,
  setToken,
  getUserInfo,
  setUserInfo,
  clearAuth,
  isAuthenticated
}
