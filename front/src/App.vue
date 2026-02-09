<script setup>
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { RouterView, useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const mobileMenuOpen = ref(false)

const toggleMobileMenu = () => {
  mobileMenuOpen.value = !mobileMenuOpen.value
}

const closeMobileMenu = () => {
  mobileMenuOpen.value = false
}

// 监听路由变化，关闭移动菜单
watch(() => route.path, () => {
  closeMobileMenu()
})

// 监听窗口大小变化
const handleResize = () => {
  if (window.innerWidth > 768) {
    closeMobileMenu()
  }
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})

// 退出登录
const handleLogout = () => {
  if (confirm('确定要退出登录吗？')) {
    userStore.logout()
    router.push('/login')
    closeMobileMenu()
  }
}

// 跳转到登录页
const goToLogin = () => {
  router.push('/login')
  closeMobileMenu()
}
</script>

<template>
  <div id="app">
    <header class="app-header">
      <div class="header-content">
        <h1 class="logo">🗑️ 垃圾桶定位系统</h1>
        <nav class="nav desktop-nav">
          <router-link to="/" class="nav-link">首页</router-link>
          <router-link to="/upload" class="nav-link">上传位置</router-link>
          <!-- 用户信息 -->
          <div v-if="userStore.isAuthenticated" class="user-info-nav">
            <router-link to="/profile" class="nav-link">用户中心</router-link>
            <span class="username-nav">{{ userStore.userInfo?.username }}</span>
            <button @click="handleLogout" class="btn-logout-nav">退出</button>
          </div>
          <div v-else class="user-info-nav">
            <button @click="goToLogin" class="btn-login-nav">登录</button>
          </div>
        </nav>
        <button class="mobile-menu-btn" @click="toggleMobileMenu" aria-label="菜单">
          <span class="hamburger-icon" :class="{ active: mobileMenuOpen }">
            <span></span>
            <span></span>
            <span></span>
          </span>
        </button>
      </div>
      <nav class="nav mobile-nav" :class="{ open: mobileMenuOpen }">
        <router-link to="/" class="nav-link" @click="closeMobileMenu">首页</router-link>
        <router-link to="/upload" class="nav-link" @click="closeMobileMenu">上传位置</router-link>
        <router-link v-if="userStore.isAuthenticated" to="/profile" class="nav-link" @click="closeMobileMenu">用户中心</router-link>
        <!-- 移动端用户信息 -->
        <div v-if="userStore.isAuthenticated" class="user-info-nav mobile-user-info">
          <span class="username-nav">{{ userStore.userInfo?.username }}</span>
          <button @click="handleLogout" class="btn-logout-nav">退出</button>
        </div>
        <div v-else class="user-info-nav mobile-user-info">
          <button @click="goToLogin" class="btn-login-nav">登录</button>
        </div>
      </nav>
    </header>
    <main class="app-main">
      <RouterView />
    </main>
    <div v-if="mobileMenuOpen" class="mobile-menu-overlay" @click="closeMobileMenu"></div>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html {
  /* 移动端优化 */
  -webkit-text-size-adjust: 100%;
  -ms-text-size-adjust: 100%;
  /* 平滑滚动 */
  scroll-behavior: smooth;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background-color: var(--bg-primary);
  color: var(--text-primary);
  /* 防止移动端双击缩放 */
  touch-action: manipulation;
  /* 移动端优化 */
  -webkit-tap-highlight-color: transparent;
}

#app {
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.app-header {
  background: var(--color-primary);
  color: white;
  padding: 0 20px;
  box-shadow: var(--shadow-md);
  z-index: 1000;
  position: relative;
}

.header-content {
  max-width: 100%;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 60px;
}

.logo {
  font-size: 20px;
  font-weight: 600;
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.nav {
  display: flex;
  gap: 20px;
}

.nav-link {
  color: white;
  text-decoration: none;
  padding: 8px 16px;
  border-radius: 4px;
  transition: background 0.3s;
  font-size: 14px;
  display: inline-block;
  /* 确保触摸目标足够大 */
  min-height: 44px;
  display: flex;
  align-items: center;
}

.nav-link:hover {
  background: rgba(255, 255, 255, 0.2);
}

.nav-link.router-link-active {
  background: rgba(255, 255, 255, 0.3);
  font-weight: 500;
}

.app-main {
  flex: 1;
  overflow: hidden;
  position: relative;
}

/* 移动端菜单按钮 */
.mobile-menu-btn {
  display: none;
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
  color: white;
  /* 确保触摸目标足够大 */
  min-width: 44px;
  min-height: 44px;
  align-items: center;
  justify-content: center;
}

.hamburger-icon {
  display: flex;
  flex-direction: column;
  gap: 5px;
  width: 24px;
  height: 18px;
}

.hamburger-icon span {
  display: block;
  width: 100%;
  height: 2px;
  background: white;
  border-radius: 2px;
  transition: all 0.3s ease;
}

.hamburger-icon.active span:nth-child(1) {
  transform: rotate(45deg) translate(7px, 7px);
}

.hamburger-icon.active span:nth-child(2) {
  opacity: 0;
}

.hamburger-icon.active span:nth-child(3) {
  transform: rotate(-45deg) translate(7px, -7px);
}

/* 移动端导航菜单 */
.mobile-nav {
  display: none;
  flex-direction: column;
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: var(--color-primary);
  box-shadow: var(--shadow-lg);
  max-height: 0;
  overflow: hidden;
  transition: max-height 0.3s ease;
}

.mobile-nav.open {
  max-height: 200px;
}

.mobile-nav .nav-link {
  width: 100%;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

/* 用户信息导航样式 */
.user-info-nav {
  display: flex;
  align-items: center;
  gap: 12px;
}

.username-nav {
  color: white;
  font-size: 14px;
  font-weight: 500;
}

.btn-logout-nav,
.btn-login-nav {
  padding: 6px 12px;
  border: 1px solid white;
  background: transparent;
  color: white;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background 0.3s;
  min-height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-logout-nav:hover,
.btn-login-nav:hover {
  background: rgba(255, 255, 255, 0.2);
}

.mobile-user-info {
  width: 100%;
  padding: 16px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  justify-content: space-between;
}

.mobile-menu-overlay {
  display: none;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 999;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .desktop-nav {
    display: none;
  }

  .mobile-menu-btn {
    display: flex;
  }

  .mobile-nav {
    display: flex;
  }

  .mobile-menu-overlay {
    display: block;
  }

  .logo {
    font-size: 18px;
  }

  .header-content {
    padding: 0 10px;
  }
}

@media (max-width: 480px) {
  .logo {
    font-size: 16px;
  }
}

/* 全局移动端优化 */
@media (max-width: 768px) {
  /* 确保所有按钮和可点击元素有足够的触摸目标 */
  button,
  a,
  [role="button"] {
    min-height: 44px;
    min-width: 44px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
  }

  /* 优化输入框在移动端的显示 */
  input,
  textarea,
  select {
    font-size: 16px; /* 防止iOS自动缩放 */
  }

  /* 优化滚动条 */
  ::-webkit-scrollbar {
    width: 6px;
    height: 6px;
  }

  ::-webkit-scrollbar-thumb {
    background: rgba(139, 111, 71, 0.3);
    border-radius: 3px;
  }

  ::-webkit-scrollbar-thumb:hover {
    background: rgba(139, 111, 71, 0.5);
  }
}

/* 防止移动端长按选择 */
@media (max-width: 768px) {
  img {
    -webkit-user-select: none;
    user-select: none;
    -webkit-touch-callout: none;
  }
}

/* 优化移动端焦点样式 */
@media (max-width: 768px) {
  *:focus {
    outline: none;
  }

  button:focus-visible,
  a:focus-visible,
  input:focus-visible,
  textarea:focus-visible {
    outline: 2px solid var(--color-primary);
    outline-offset: 2px;
  }
}
</style>
