<template>
  <div class="profile-container">
    <div class="profile-card">
      <h1 class="profile-title">用户中心</h1>
      
      <!-- 用户信息展示 -->
      <div class="user-info-section">
        <div class="user-info-item">
          <span class="label">用户名：</span>
          <span class="value">{{ userStore.userInfo?.username || '未知' }}</span>
        </div>
        <div class="user-info-item" v-if="userStore.userInfo?.created_at">
          <span class="label">注册时间：</span>
          <span class="value">{{ formatDate(userStore.userInfo.created_at) }}</span>
        </div>
      </div>

      <!-- 标签页导航 -->
      <div class="tab-buttons">
        <button 
          :class="['tab-btn', { active: activeTab === 'my-trashcans' }]"
          @click="activeTab = 'my-trashcans'"
        >
          我的上传
        </button>
      </div>

      <!-- 标签页内容 -->
      <div class="tab-content">
        <MyTrashCans v-if="activeTab === 'my-trashcans'" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import MyTrashCans from '@/views/MyTrashCans.vue'

const router = useRouter()
const userStore = useUserStore()
const activeTab = ref('my-trashcans')

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  // 如果未登录，跳转到登录页
  if (!userStore.isAuthenticated) {
    router.push('/login')
  } else {
    // 确保用户信息已加载
    if (!userStore.userInfo) {
      userStore.initUser()
    }
  }
})
</script>

<style scoped>
.profile-container {
  min-height: calc(100vh - 60px);
  padding: 20px;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  background: var(--bg-primary);
}

.profile-card {
  width: 100%;
  max-width: 1200px;
  background: var(--bg-secondary);
  border-radius: 8px;
  padding: 30px;
  box-shadow: var(--shadow-md);
}

.profile-title {
  margin: 0 0 30px 0;
  font-size: 24px;
  color: var(--text-primary);
  text-align: center;
}

.user-info-section {
  margin-bottom: 30px;
  padding: 20px;
  background: var(--bg-primary);
  border-radius: 8px;
  border: 1px solid var(--border-color);
}

.user-info-item {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
  font-size: 16px;
}

.user-info-item:last-child {
  margin-bottom: 0;
}

.user-info-item .label {
  color: var(--text-secondary);
  font-weight: 500;
  min-width: 100px;
}

.user-info-item .value {
  color: var(--text-primary);
}

.tab-buttons {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
  border-bottom: 2px solid var(--border-color);
}

.tab-btn {
  padding: 12px 24px;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  font-size: 16px;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  margin-bottom: -2px;
  transition: var(--transition-base);
}

.tab-btn:hover {
  color: var(--text-primary);
}

.tab-btn.active {
  color: var(--color-primary);
  border-bottom-color: var(--color-primary);
  font-weight: 500;
}

.tab-content {
  min-height: 400px;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .profile-container {
    padding: 10px;
  }

  .profile-card {
    padding: 20px;
  }

  .profile-title {
    font-size: 20px;
  }

  .user-info-section {
    padding: 15px;
  }

  .user-info-item {
    font-size: 14px;
    flex-direction: column;
    align-items: flex-start;
  }

  .user-info-item .label {
    min-width: auto;
    margin-bottom: 5px;
  }

  .tab-btn {
    padding: 10px 16px;
    font-size: 14px;
  }
}
</style>

