const { request } = require("../request")

function register(username, password) {
  return request({
    url: "/users/register",
    method: "POST",
    data: { username, password }
  })
}

function login(username, password) {
  return request({
    url: "/users/login",
    method: "POST",
    data: { username, password }
  })
}

function getCurrentUser() {
  return request({
    url: "/users/me",
    method: "GET"
  })
}

module.exports = {
  register,
  login,
  getCurrentUser
}
