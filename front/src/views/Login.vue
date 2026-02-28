<template>
  <div class="login-container">
    <div class="login-card">
      <h1 class="login-title">垃圾桶定位系统</h1>
      
      <div class="tab-buttons">
        <button 
          :class="['tab-btn', { active: mode === 'login' }]"
          @click="mode = 'login'"
        >
          登录
        </button>
        <button 
          :class="['tab-btn', { active: mode === 'register' }]"
          @click="mode = 'register'"
        >
          注册
        </button>
      </div>

      <!-- 登录表单 -->
      <form v-if="mode === 'login'" @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label>用户名</label>
          <input 
            v-model="loginForm.username" 
            type="text" 
            placeholder="请输入用户名"
            class="input"
            required
          />
        </div>
        <div class="form-group">
          <label>密码</label>
          <input 
            v-model="loginForm.password" 
            type="password" 
            placeholder="请输入密码"
            class="input"
            required
          />
        </div>
        <button type="submit" :disabled="loading" class="btn btn-primary">
          {{ loading ? '登录中...' : '登录' }}
        </button>
        <p v-if="error" class="error-message">{{ error }}</p>
      </form>

      <!-- 注册表单 -->
      <form v-if="mode === 'register'" @submit.prevent="handleRegister" class="login-form">
        <div class="form-group">
          <label>用户名</label>
          <input 
            v-model="registerForm.username" 
            type="text" 
            placeholder="请输入用户名（3-20个字符）"
            class="input"
            required
            minlength="3"
            maxlength="20"
          />
        </div>
        <div class="form-group">
          <label>密码</label>
          <input 
            v-model="registerForm.password" 
            type="password" 
            placeholder="请输入密码（至少6个字符）"
            class="input"
            required
            minlength="6"
          />
        </div>
        <div class="form-group">
          <label>确认密码</label>
          <input 
            v-model="registerForm.confirmPassword" 
            type="password" 
            placeholder="请再次输入密码"
            class="input"
            required
          />
        </div>
        <button type="submit" :disabled="loading" class="btn btn-primary">
          {{ loading ? '注册中...' : '注册' }}
        </button>
        <p v-if="error" class="error-message">{{ error }}</p>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const mode = ref('login')
const loading = ref(false)
const error = ref('')

const loginForm = reactive({
  username: '',
  password: ''
})

const registerForm = reactive({
  username: '',
  password: '',
  confirmPassword: ''
})

const handleLogin = async () => {
  error.value = ''
  loading.value = true
  
  try {
    const result = await userStore.loginUser(loginForm.username, loginForm.password)
    if (result.success) {
      // 登录成功，跳转到首页
      router.push('/')
    } else {
      error.value = result.message
    }
  } catch (err) {
    error.value = err.message || '登录失败'
  } finally {
    loading.value = false
  }
}

const handleRegister = async () => {
  error.value = ''
  
  // 验证密码确认
  if (registerForm.password !== registerForm.confirmPassword) {
    error.value = '两次输入的密码不一致'
    return
  }
  
  loading.value = true
  
  try {
    const result = await userStore.registerUser(registerForm.username, registerForm.password)
    if (result.success) {
      // 注册成功，跳转到首页
      router.push('/')
    } else {
      error.value = result.message
    }
  } catch (err) {
    error.value = err.message || '注册失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: var(--bg-primary);
  padding: 20px;
  position: relative;
  z-index: 1;
  overflow: auto;
}

.login-card {
  width: 100%;
  max-width: 400px;
  background: var(--bg-secondary);
  border-radius: 8px;
  padding: 30px;
  box-shadow: var(--shadow-lg);
  border: 1px solid var(--border-color);
  position: relative;
  z-index: 1;
}

.login-title {
  text-align: center;
  margin: 0 0 30px 0;
  font-size: 24px;
  color: var(--text-primary);
}

.tab-buttons {
  display: flex;
  gap: 10px;
  margin-bottom: 30px;
}

.tab-btn {
  flex: 1;
  padding: 10px;
  border: 1px solid var(--border-color);
  background: var(--bg-primary);
  color: var(--text-secondary);
  border-radius: 4px;
  cursor: pointer;
  transition: var(--transition-base);
  font-size: 14px;
}

.tab-btn:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.tab-btn.active {
  background: var(--color-primary);
  color: white;
  border-color: var(--color-primary);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
  position: relative;
  z-index: 1;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
  position: relative;
  z-index: 1;
}

.form-group label {
  color: var(--text-primary);
  font-size: 14px;
  font-weight: 500;
}

.input {
  width: 100%;
  padding: 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
  background: var(--bg-primary);
  color: var(--text-primary);
  transition: var(--transition-base);
  pointer-events: auto;
  -webkit-user-select: text;
  user-select: text;
  position: relative;
  z-index: 1;
}

.input:focus {
  outline: none;
  border-color: var(--border-color-focus);
  box-shadow: 0 0 0 2px rgba(139, 111, 71, 0.1);
}

.btn {
  padding: 12px 20px;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: var(--transition-base);
  font-weight: 500;
  min-height: 44px;
}

.btn-primary {
  background: var(--color-primary);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: var(--color-primary-light);
  box-shadow: var(--shadow-sm);
}

.btn-primary:disabled {
  background: var(--text-disabled);
  cursor: not-allowed;
}

.error-message {
  color: var(--color-error);
  font-size: 14px;
  margin: 0;
  text-align: center;
}

@media (max-width: 480px) {
  .login-card {
    padding: 20px;
  }
  
  .login-title {
    font-size: 20px;
  }
}
</style>

