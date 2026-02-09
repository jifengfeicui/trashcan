import axios from 'axios'

// 创建axios实例
const request = axios.create({
  baseURL: '/api', // 通过Vite代理转发到后端
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    // 添加token到请求头
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => {
    // 根据后端返回的数据结构处理
    // 后端统一返回格式为 { code, data, msg }
    // SUCCESS = 2000
    if (response.data && response.data.code !== undefined) {
      if (response.data.code === 2000) {
        return response.data
      } else {
        return Promise.reject(new Error(response.data.msg || '请求失败'))
      }
    }
    return response.data
  },
  error => {
    // 处理401未授权错误
    if (error.response && error.response.status === 401) {
      // 清除token
      localStorage.removeItem('token')
      // 跳转到登录页（如果不在登录页）
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }
    // 处理错误响应
    const message = error.response?.data?.msg || error.message || '请求失败'
    console.error('API请求错误:', message)
    return Promise.reject(error)
  }
)

export default request

