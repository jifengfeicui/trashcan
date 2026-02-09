import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'
import Upload from '@/views/Upload.vue'
import Login from '@/views/Login.vue'
import Profile from '@/views/Profile.vue'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/upload',
      name: 'upload',
      component: Upload,
      meta: { requiresAuth: true }
    },
    {
      path: '/login',
      name: 'login',
      component: Login
    },
    {
      path: '/profile',
      name: 'profile',
      component: Profile,
      meta: { requiresAuth: true }
    }
  ],
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  
  // 如果路由需要认证
  if (to.meta.requiresAuth) {
    // 检查是否已登录
    if (!userStore.isAuthenticated) {
      // 未登录，跳转到登录页
      next({ name: 'login', query: { redirect: to.fullPath } })
    } else {
      next()
    }
  } else {
    // 如果已登录且访问登录页，跳转到首页
    if (to.name === 'login' && userStore.isAuthenticated) {
      next('/')
    } else {
      next()
    }
  }
})

export default router
