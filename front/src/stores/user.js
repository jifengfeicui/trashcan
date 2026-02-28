import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { register, login, getCurrentUser } from '@/api/user'

export const useUserStore = defineStore('user', () => {
  // 状态
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(null)

  // 计算属性
  const isAuthenticated = computed(() => !!token.value)

  // 初始化时尝试获取用户信息
  const initUser = async () => {
    if (token.value) {
      try {
        const response = await getCurrentUser()
        if (response.code === 2000 && response.data) {
          userInfo.value = response.data
        } else {
          // token无效，清除
          logout()
        }
      } catch (error) {
        // token无效，清除
        logout()
      }
    }
  }

  // 注册
  const registerUser = async (username, password) => {
    try {
      const response = await register(username, password)
      if (response.code === 2000 && response.data) {
        token.value = response.data.token
        userInfo.value = response.data.user
        localStorage.setItem('token', response.data.token)
        return { success: true, message: response.msg || '注册成功' }
      } else {
        return { success: false, message: response.msg || '注册失败' }
      }
    } catch (error) {
      return { success: false, message: error.message || '注册失败' }
    }
  }

  // 登录
  const loginUser = async (username, password) => {
    try {
      const response = await login(username, password)
      if (response.code === 2000 && response.data) {
        token.value = response.data.token
        userInfo.value = response.data.user
        localStorage.setItem('token', response.data.token)
        return { success: true, message: response.msg || '登录成功' }
      } else {
        return { success: false, message: response.msg || '登录失败' }
      }
    } catch (error) {
      return { success: false, message: error.message || '登录失败' }
    }
  }

  // 退出登录
  const logout = () => {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
  }

  // 检查认证状态
  const checkAuth = async () => {
    if (!token.value) {
      return false
    }
    try {
      const response = await getCurrentUser()
      if (response.code === 2000 && response.data) {
        userInfo.value = response.data
        return true
      } else {
        logout()
        return false
      }
    } catch (error) {
      logout()
      return false
    }
  }

  return {
    token,
    userInfo,
    isAuthenticated,
    registerUser,
    loginUser,
    logout,
    checkAuth,
    initUser
  }
})

