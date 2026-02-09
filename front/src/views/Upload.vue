<template>
  <div class="upload-container">
    <!-- 移动端侧边栏切换按钮 -->
    <button class="mobile-sidebar-toggle" @click="toggleSidebar" aria-label="切换侧边栏">
      <span class="toggle-icon" :class="{ active: sidebarOpen }">
        <span></span>
        <span></span>
        <span></span>
      </span>
    </button>

    <!-- 侧边栏遮罩层 -->
    <div v-if="sidebarOpen" class="sidebar-overlay" @click="closeSidebar"></div>

    <div class="upload-form-panel" :class="{ open: sidebarOpen }">
      <div class="upload-form-panel-header">
        <h2>上传垃圾桶位置</h2>
        <button class="close-sidebar-btn" @click="closeSidebar" aria-label="关闭">
          <span>×</span>
        </button>
      </div>
      <div class="upload-form-panel-content">
        <form @submit.prevent="handleSubmit" class="upload-form">
        <div class="form-group">
          <label>位置选择方式:</label>
          <div class="radio-group">
            <label>
              <input type="radio" v-model="locationMode" value="current" />
              当前定位
            </label>
            <label>
              <input type="radio" v-model="locationMode" value="manual" />
              手动输入
            </label>
          </div>
        </div>

        <div v-if="locationMode === 'current'" class="form-group">
          <button 
            type="button" 
            @click="getCurrentLocation" 
            :disabled="locating"
            class="btn btn-location"
          >
            {{ locating ? '定位中...' : '📍 获取当前位置' }}
          </button>
          <p v-if="currentLocation" class="selected-location">
            当前位置: {{ currentLocation.lng.toFixed(6) }}, {{ currentLocation.lat.toFixed(6) }}
          </p>
          <p v-else class="hint">点击按钮获取您的当前位置</p>
        </div>

        <div v-if="locationMode === 'manual'" class="form-group">
          <label>经度:</label>
          <input 
            v-model.number="formData.longitude" 
            type="number" 
            step="any"
            placeholder="例如: 116.397428"
            class="input"
            required
          />
        </div>

        <div v-if="locationMode === 'manual'" class="form-group">
          <label>纬度:</label>
          <input 
            v-model.number="formData.latitude" 
            type="number" 
            step="any"
            placeholder="例如: 39.90923"
            class="input"
            required
          />
        </div>

        <div class="form-group">
          <label>地址描述:</label>
          <input 
            v-model="formData.address" 
            type="text" 
            placeholder="例如: 北京市朝阳区xxx街道"
            class="input"
          />
        </div>

        <div class="form-group">
          <label>详细描述:</label>
          <textarea 
            v-model="formData.description" 
            rows="3"
            placeholder="可选：垃圾桶的详细描述信息"
            class="textarea"
          ></textarea>
        </div>

        <div class="form-group">
          <label>上传图片:</label>
          <div class="upload-area" @click="triggerFileInput">
            <input 
              ref="fileInput"
              type="file" 
              accept="image/*"
              @change="handleFileChange"
              style="display: none"
            />
            <div v-if="!imagePreview" class="upload-placeholder">
              <p>点击选择图片</p>
              <p class="hint-text">支持 JPG、PNG 等格式</p>
            </div>
            <div v-else class="image-preview">
              <img :src="imagePreview" alt="预览图片" />
              <button type="button" @click.stop="removeImage" class="remove-btn">×</button>
            </div>
          </div>
        </div>

        <div class="form-actions">
          <button type="submit" :disabled="submitting" class="btn btn-primary">
            {{ submitting ? '提交中...' : '提交' }}
          </button>
          <button type="button" @click="resetForm" class="btn btn-secondary">
            重置
          </button>
        </div>
      </form>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, reactive, watch, onMounted, onUnmounted } from 'vue'
import { createTrashCan } from '@/api/trashcan'

const fileInput = ref(null)
const locationMode = ref('current')
const currentLocation = ref(null)
const locating = ref(false)
const imagePreview = ref(null)
const selectedFile = ref(null)
const submitting = ref(false)
const sidebarOpen = ref(false)

const toggleSidebar = () => {
  sidebarOpen.value = !sidebarOpen.value
}

const closeSidebar = () => {
  sidebarOpen.value = false
}

// 监听窗口大小变化
const handleResize = () => {
  if (window.innerWidth > 768) {
    sidebarOpen.value = false
  } else {
    // 移动端时，操作面板始终显示（通过CSS控制），不需要侧边栏切换
    sidebarOpen.value = true
  }
}

onMounted(() => {
  window.addEventListener('resize', handleResize)
  // 初始化时检查是否为移动端
  if (window.innerWidth <= 768) {
    sidebarOpen.value = true
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})

const formData = reactive({
  latitude: null,
  longitude: null,
  address: '',
  description: ''
})

// 获取当前位置
const getCurrentLocation = () => {
  if (!navigator.geolocation) {
    alert('您的浏览器不支持地理定位功能')
    return
  }

  locating.value = true

  navigator.geolocation.getCurrentPosition(
    (position) => {
      const lat = position.coords.latitude
      const lng = position.coords.longitude
      const location = { lng, lat }
      currentLocation.value = location
      formData.latitude = lat
      formData.longitude = lng
      locating.value = false
    },
    (error) => {
      locating.value = false
      let errorMessage = '定位失败'
      switch (error.code) {
        case error.PERMISSION_DENIED:
          errorMessage = '定位权限被拒绝，请在浏览器设置中允许定位权限'
          break
        case error.POSITION_UNAVAILABLE:
          errorMessage = '定位信息不可用'
          break
        case error.TIMEOUT:
          errorMessage = '定位请求超时'
          break
      }
      alert(errorMessage)
      console.error('定位错误:', error)
    },
    {
      enableHighAccuracy: true,
      timeout: 10000,
      maximumAge: 0
    }
  )
}

// 监听位置模式变化
watch(locationMode, (newMode) => {
  if (newMode === 'current' && currentLocation.value) {
    formData.latitude = currentLocation.value.lat
    formData.longitude = currentLocation.value.lng
  }
})

// 触发文件选择
const triggerFileInput = () => {
  fileInput.value?.click()
}

// 处理文件选择
const handleFileChange = (e) => {
  const file = e.target.files[0]
  if (file) {
    // 验证文件类型
    if (!file.type.startsWith('image/')) {
      alert('请选择图片文件')
      return
    }
    
    // 验证文件大小（限制5MB）
    if (file.size > 5 * 1024 * 1024) {
      alert('图片大小不能超过5MB')
      return
    }
    
    selectedFile.value = file
    
    // 预览图片
    const reader = new FileReader()
    reader.onload = (e) => {
      imagePreview.value = e.target.result
    }
    reader.readAsDataURL(file)
  }
}

// 移除图片
const removeImage = () => {
  imagePreview.value = null
  selectedFile.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

// 重置表单
const resetForm = () => {
  formData.latitude = null
  formData.longitude = null
  formData.address = ''
  formData.description = ''
  currentLocation.value = null
  locating.value = false
  removeImage()
}

// 提交表单
const handleSubmit = async () => {
  // 验证必填字段
  if (!formData.latitude || !formData.longitude) {
    alert('请选择或输入位置信息')
    return
  }

  submitting.value = true
  
  try {
    // 创建FormData
    const formDataToSend = new FormData()
    formDataToSend.append('latitude', formData.latitude.toString())
    formDataToSend.append('longitude', formData.longitude.toString())
    if (formData.address) {
      formDataToSend.append('address', formData.address)
    }
    if (formData.description) {
      formDataToSend.append('description', formData.description)
    }
    if (selectedFile.value) {
      formDataToSend.append('image', selectedFile.value)
    }

    const response = await createTrashCan(formDataToSend)
    
    if (response.code === 2000) {
      alert('上传成功！')
      resetForm()
    } else {
      alert(response.msg || '上传失败')
    }
  } catch (error) {
    console.error('上传错误:', error)
    alert('上传失败: ' + error.message)
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.upload-container {
  display: flex;
  height: 100%;
  overflow-y: auto;
  position: relative;
  justify-content: center;
  align-items: flex-start;
  padding: 0 !important;
  margin: 0 !important;
  margin-top: 0 !important;
  width: 100%;
  min-height: 0;
  top: 0;
}

/* 移动端侧边栏切换按钮 */
.mobile-sidebar-toggle {
  display: none;
  position: fixed;
  top: 70px;
  left: 10px;
  z-index: 1001;
  background: var(--color-primary);
  color: white;
  border: none;
  border-radius: 8px;
  padding: 10px;
  cursor: pointer;
  box-shadow: var(--shadow-lg);
  min-width: 44px;
  min-height: 44px;
  align-items: center;
  justify-content: center;
}

.toggle-icon {
  display: flex;
  flex-direction: column;
  gap: 4px;
  width: 20px;
  height: 16px;
}

.toggle-icon span {
  display: block;
  width: 100%;
  height: 2px;
  background: white;
  border-radius: 2px;
  transition: all 0.3s ease;
}

.toggle-icon.active span:nth-child(1) {
  transform: rotate(45deg) translate(6px, 6px);
}

.toggle-icon.active span:nth-child(2) {
  opacity: 0;
}

.toggle-icon.active span:nth-child(3) {
  transform: rotate(-45deg) translate(6px, -6px);
}

/* 侧边栏遮罩层 */
.sidebar-overlay {
  display: none;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 999;
}

.upload-form-panel {
  width: 100%;
  max-width: 600px;
  background: var(--bg-primary);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  transition: transform 0.3s ease;
  margin: 0;
  padding: 0;
  border-radius: 0;
  box-shadow: none;
  align-self: flex-start;
  position: relative;
  top: 0;
}

.upload-form-panel-header {
  display: flex;
  padding: 8px 20px;
  border-bottom: 1px solid var(--border-color);
  align-items: center;
  justify-content: space-between;
  background: var(--bg-secondary);
  margin: 0;
}

.upload-form-panel-header h2 {
  margin: 0;
  font-size: 18px;
  color: var(--text-primary);
}

.close-sidebar-btn {
  background: none;
  border: none;
  font-size: 28px;
  color: var(--text-secondary);
  cursor: pointer;
  padding: 0;
  width: 32px;
  height: 32px;
  display: none; /* 桌面端默认隐藏 */
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: var(--transition-base);
}

.close-sidebar-btn:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.upload-form-panel-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  /* 移动端滚动优化 */
  -webkit-overflow-scrolling: touch;
}

.upload-form-panel h2 {
  margin: 0 0 20px 0;
  font-size: 20px;
  color: var(--text-primary);
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: var(--text-primary);
  font-size: 14px;
  font-weight: 500;
}

.radio-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.radio-group label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: normal;
  cursor: pointer;
  padding: 8px;
  border-radius: 4px;
  transition: var(--transition-base);
  /* 确保触摸目标足够大 */
  min-height: 44px;
}

.radio-group label:hover {
  background: var(--bg-tertiary);
}

.radio-group input[type="radio"] {
  width: 20px;
  height: 20px;
  cursor: pointer;
  accent-color: var(--color-primary);
}

.input,
.textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
  font-family: inherit;
  background: var(--bg-secondary);
  color: var(--text-primary);
  transition: var(--transition-base);
  /* 移动端优化 */
  -webkit-appearance: none;
  appearance: none;
}

.textarea {
  resize: vertical;
}

.input:focus,
.textarea:focus {
  outline: none;
  border-color: var(--border-color-focus);
  box-shadow: 0 0 0 2px rgba(139, 111, 71, 0.1);
}

.hint {
  color: var(--text-secondary);
  font-size: 13px;
  margin: 5px 0;
}

.selected-location {
  color: var(--color-primary);
  font-size: 13px;
  margin: 5px 0;
  font-weight: 500;
}

.upload-area {
  border: 2px dashed var(--border-color);
  border-radius: 4px;
  padding: 20px;
  text-align: center;
  cursor: pointer;
  transition: var(--transition-base);
  background: var(--bg-secondary);
}

.upload-area:hover {
  border-color: var(--border-color-hover);
  background: var(--bg-tertiary);
}

.upload-placeholder p {
  margin: 5px 0;
  color: var(--text-secondary);
}

.hint-text {
  font-size: 12px;
  color: var(--text-tertiary);
}

.image-preview {
  position: relative;
  display: inline-block;
}

.image-preview img {
  max-width: 100%;
  max-height: 300px;
  border-radius: 4px;
}

.remove-btn {
  position: absolute;
  top: -10px;
  right: -10px;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--color-error);
  color: white;
  border: none;
  cursor: pointer;
  font-size: 18px;
  line-height: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: var(--transition-base);
}

.remove-btn:hover {
  background: var(--color-error-light);
}

.form-actions {
  display: flex;
  gap: 10px;
  margin-top: 30px;
}

.btn {
  flex: 1;
  padding: 12px 20px;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
  /* 确保触摸目标足够大 */
  min-height: 44px;
  font-weight: 500;
}

.btn-primary {
  background: var(--color-primary);
  color: white;
  transition: var(--transition-base);
}

.btn-primary:hover:not(:disabled) {
  background: var(--color-primary-light);
  box-shadow: var(--shadow-sm);
}

.btn-primary:active:not(:disabled) {
  transform: scale(0.98);
}

.btn-primary:disabled {
  background: var(--text-disabled);
  cursor: not-allowed;
}

.btn-secondary {
  background: var(--bg-secondary);
  color: var(--color-primary);
  border: 1px solid var(--color-primary);
  transition: var(--transition-base);
}

.btn-secondary:hover {
  background: var(--bg-tertiary);
  border-color: var(--color-primary-light);
}

.btn-secondary:active {
  transform: scale(0.98);
}

.btn-location {
  width: 100%;
  background: var(--color-success);
  color: white;
  border: none;
  padding: 12px 20px;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: var(--transition-base);
}

.btn-location:hover:not(:disabled) {
  background: var(--color-success-light);
  box-shadow: var(--shadow-sm);
}

.btn-location:active:not(:disabled) {
  transform: scale(0.98);
}

.btn-location:disabled {
  background: var(--text-disabled);
  cursor: not-allowed;
}

.upload-area {
  /* 移动端触摸优化 */
  -webkit-tap-highlight-color: transparent;
}

.upload-area:active {
  transform: scale(0.98);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .upload-container {
    flex-direction: column;
    height: 100vh;
    overflow-y: auto;
  }

  .mobile-sidebar-toggle {
    display: none; /* 移动端不再需要侧边栏切换按钮 */
  }

  .sidebar-overlay {
    display: none; /* 移动端不再需要遮罩层 */
  }

  .upload-form-panel {
    width: 100%;
    margin: 0;
    border-radius: 0;
    box-shadow: none;
    max-width: 100%;
  }

  .upload-form-panel.open {
    transform: none;
  }

  .upload-form-panel-header {
    display: flex;
  }

  /* 移动端显示关闭按钮 */
  .close-sidebar-btn {
    display: flex;
  }

  .upload-form-panel-content {
    padding: 15px;
  }

  .radio-group {
    gap: 10px;
  }
}

@media (max-width: 480px) {
  .upload-form-panel {
    width: 100%;
  }

  .upload-form-panel-content {
    padding: 12px;
  }

  .form-group {
    margin-bottom: 15px;
  }
}
</style>

