<template>
  <div class="home-container">
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

    <div class="search-panel" :class="{ open: sidebarOpen }">
      <div class="search-panel-header">
        <h2>查找附近垃圾桶</h2>
        <button class="close-sidebar-btn" @click="closeSidebar" aria-label="关闭">
          <span>×</span>
        </button>
      </div>
      <div class="search-panel-content">
        <div class="search-form">
          <div class="form-row">
            <div class="form-group form-group-inline">
              <label>搜索半径（公里）:</label>
              <input 
                v-model.number="searchRadius" 
                type="number" 
                min="1" 
                max="50" 
                step="1"
                class="input"
              />
            </div>
            <div class="form-group form-group-inline">
              <label>返回数量:</label>
              <input 
                v-model.number="searchLimit" 
                type="number" 
                min="1" 
                max="50" 
                step="1"
                class="input"
              />
            </div>
          </div>
          <div class="button-group button-group-row">
            <button @click="searchNearby" :disabled="loading" class="btn btn-primary">
              {{ loading ? '搜索中...' : '搜索附近垃圾桶' }}
            </button>
            <button @click="locateMe" class="btn btn-secondary">
              定位到我的位置
            </button>
          </div>
        </div>
        
        <!-- 桌面端显示搜索结果列表 -->
        <div v-if="trashCans.length > 0 && !isMobile" class="results-list">
          <h3>搜索结果 ({{ trashCans.length }})</h3>
          <div 
            v-for="item in trashCans" 
            :key="item.id" 
            class="result-item"
            @click="focusMarker(item)"
          >
            <div class="result-content">
              <h4>{{ item.address || '垃圾桶位置' }}</h4>
              <p v-if="item.description">{{ item.description }}</p>
              <p class="distance">距离: {{ item.distance?.toFixed(2) || '未知' }} 公里</p>
            </div>
            <img 
              v-if="item.image_url" 
              :src="item.image_url" 
              alt="垃圾桶图片"
              class="result-image"
            />
          </div>
        </div>
        
        <!-- 手机端显示搜索结果按钮 -->
        <div v-if="trashCans.length > 0 && isMobile" class="mobile-results-button">
          <button @click="openResultsModal" class="btn btn-primary">
            查看搜索结果 ({{ trashCans.length }})
          </button>
        </div>
      </div>
    </div>
    
    <div class="map-panel">
      <MapContainer 
        ref="mapRef"
        :center="mapCenter" 
        :markers="trashCans"
        :zoom="15"
        @map-ready="onMapReady"
        @location-ready="onLocationReady"
        @map-click="onMapClick"
        @info-window-action="onInfoWindowAction"
        @navigation-ready="onNavigationReady"
      />
      
      <!-- 导航信息面板 -->
      <div v-if="navigationInfo" class="navigation-panel">
        <div class="navigation-header">
          <h4>步行路线</h4>
          <button @click="closeNavigation" class="close-nav-btn" aria-label="关闭导航">×</button>
        </div>
        <div class="navigation-content">
          <div class="nav-info-item">
            <span class="nav-label">距离:</span>
            <span class="nav-value">{{ navigationInfo.distance }} 公里</span>
          </div>
          <div class="nav-info-item">
            <span class="nav-label">预计时间:</span>
            <span class="nav-value">{{ navigationInfo.duration }} 分钟</span>
          </div>
          <div class="nav-actions">
            <button @click="openExternalNavigation" class="btn btn-primary btn-sm">
              打开高德地图导航
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 手机端搜索结果弹框 -->
    <div v-if="isMobile" class="mobile-results-modal" :class="{ open: resultsModalOpen }">
      <div class="modal-overlay" @click="closeResultsModal"></div>
      <div class="modal-content">
        <div class="modal-header">
          <h3>搜索结果 ({{ trashCans.length }})</h3>
          <button class="modal-close-btn" @click="closeResultsModal" aria-label="关闭">
            <span>×</span>
          </button>
        </div>
        <div class="modal-body">
          <div 
            v-for="item in trashCans" 
            :key="item.id" 
            class="result-item"
            @click="focusMarkerAndClose(item)"
          >
            <div class="result-content">
              <h4>{{ item.address || '垃圾桶位置' }}</h4>
              <p v-if="item.description">{{ item.description }}</p>
              <p class="distance">距离: {{ item.distance?.toFixed(2) || '未知' }} 公里</p>
            </div>
            <img 
              v-if="item.image_url" 
              :src="item.image_url" 
              alt="垃圾桶图片"
              class="result-image"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import MapContainer from '@/components/MapContainer.vue'
import { getNearbyTrashCans } from '@/api/trashcan'

const mapRef = ref(null)
const mapCenter = ref([116.397428, 39.90923]) // 默认北京
const userLocation = ref(null)
const trashCans = ref([])
const loading = ref(false)
const searchRadius = ref(5)
const searchLimit = ref(10)
const sidebarOpen = ref(false)
const resultsModalOpen = ref(false)
const isMobile = ref(false)
const navigationInfo = ref(null)
const navigationDestination = ref(null)

const toggleSidebar = () => {
  sidebarOpen.value = !sidebarOpen.value
}

const closeSidebar = () => {
  sidebarOpen.value = false
}

// 检测是否为移动端
const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
}

// 监听窗口大小变化
const handleResize = () => {
  checkMobile()
  if (window.innerWidth > 768) {
    sidebarOpen.value = false
    resultsModalOpen.value = false // 桌面端关闭弹框
  } else {
    // 移动端时，操作面板始终显示（通过CSS控制），不需要侧边栏切换
    sidebarOpen.value = true
  }
}

// 打开搜索结果弹框
const openResultsModal = () => {
  resultsModalOpen.value = true
}

// 关闭搜索结果弹框
const closeResultsModal = () => {
  resultsModalOpen.value = false
}

// 聚焦标记并关闭弹框
const focusMarkerAndClose = (item) => {
  focusMarker(item)
  closeResultsModal()
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', handleResize)
  // 初始化时检查是否为移动端
  if (window.innerWidth <= 768) {
    sidebarOpen.value = true
  }
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})

// 地图准备就绪（仅代表地图初始化完成）
const onMapReady = () => {
  console.log('地图初始化完成')
  // 地图已准备好，可以开始使用
}

// 用户位置就绪（定位完成）
const onLocationReady = (location) => {
  userLocation.value = location
  mapCenter.value = [location.lng, location.lat]
  console.log('用户位置已获取:', location)
  // 不再自动搜索，需要用户点击按钮查询
}

// 地图点击事件
const onMapClick = (point) => {
  console.log('地图点击:', point)
}

// 搜索附近的垃圾桶
const searchNearby = async () => {
  if (!userLocation.value) {
    alert('请先获取您的位置')
    return
  }

  loading.value = true
  try {
    const response = await getNearbyTrashCans(
      userLocation.value.lat,
      userLocation.value.lng,
      searchRadius.value,
      searchLimit.value
    )
    
    if (response.code === 2000 && response.data) {
      trashCans.value = response.data
    } else {
      console.error('搜索失败:', response.msg)
      alert(response.msg || '搜索失败')
    }
  } catch (error) {
    console.error('搜索错误:', error)
    alert('搜索失败: ' + error.message)
  } finally {
    loading.value = false
  }
}

// 定位到我的位置
const locateMe = () => {
  if (navigator.geolocation) {
    navigator.geolocation.getCurrentPosition(
      (position) => {
        const { latitude, longitude } = position.coords
        userLocation.value = { lat: latitude, lng: longitude }
        mapCenter.value = [longitude, latitude]
        
        // 调用地图组件的定位方法，在地图上显示位置标记
        if (mapRef.value && mapRef.value.locateUser) {
          const locationInfo = mapRef.value.locateUser(longitude, latitude)
          console.log('定位成功，位置信息已输出到控制台')
        }
        
        // 自动搜索附近的垃圾桶
        searchNearby()
      },
      (error) => {
        console.error('定位失败:', error)
        alert('定位失败，请检查浏览器定位权限')
      },
      {
        enableHighAccuracy: true,
        timeout: 10000,
        maximumAge: 0
      }
    )
  } else {
    alert('您的浏览器不支持定位功能')
  }
}

// 聚焦到指定标记
const focusMarker = (item) => {
  if (mapRef.value && item.longitude && item.latitude) {
    mapCenter.value = [item.longitude, item.latitude]
  }
}

// 处理 InfoWindow 中的操作
const onInfoWindowAction = async ({ action, data }) => {
  if (action === 'navigate') {
    // 步行导航功能
    const { lng, lat } = data
    navigationDestination.value = { lng, lat }
    await startNavigation({ lng, lat })
  } else if (action === 'open-image') {
    // 打开图片
    const { imageUrl } = data
    window.open(imageUrl, '_blank')
  }
}

// 开始步行导航
const startNavigation = async (destination) => {
  if (!mapRef.value || !mapRef.value.navigateTo) {
    alert('地图未准备好，请稍后再试')
    return
  }

  if (!userLocation.value) {
    alert('请先获取您的位置')
    return
  }

  try {
    await mapRef.value.navigateTo(destination, userLocation.value)
  } catch (error) {
    console.error('步行路线规划失败:', error)
    
    // 如果路线规划失败，直接打开外部导航（不显示提示）
    const { lng, lat } = destination
    openExternalNavigationWithMode(lng, lat, 'walking')
  }
}

// 导航完成回调
const onNavigationReady = (info) => {
  navigationInfo.value = info
}

// 关闭导航
const closeNavigation = () => {
  if (mapRef.value && mapRef.value.clearNavigation) {
    mapRef.value.clearNavigation()
  }
  navigationInfo.value = null
  navigationDestination.value = null
}

// 打开外部导航（高德地图步行导航）
const openExternalNavigation = () => {
  if (!navigationDestination.value) return
  
  const { lng, lat } = navigationDestination.value
  openExternalNavigationWithMode(lng, lat, 'walking')
}

// 打开外部导航（步行模式）
const openExternalNavigationWithMode = (lng, lat, mode) => {
  // 高德地图App URL（Android）- 步行模式
  // 使用 route/plan 接口，添加 mode=walk 参数指定步行模式
  const appUrl = `androidamap://route/plan?sourceApplication=垃圾桶定位系统&dlat=${lat}&dlon=${lng}&dname=垃圾桶位置&mode=walk&dev=0&style=2`
  
  // 尝试打开App（不打开新网页）
  const iframe = document.createElement('iframe')
  iframe.style.display = 'none'
  iframe.src = appUrl
  document.body.appendChild(iframe)
  
  setTimeout(() => {
    document.body.removeChild(iframe)
  }, 2000)
}

onMounted(() => {
  // 尝试获取用户位置
  if (navigator.geolocation) {
    navigator.geolocation.getCurrentPosition(
      (position) => {
        const { latitude, longitude } = position.coords
        userLocation.value = { lat: latitude, lng: longitude }
        mapCenter.value = [longitude, latitude]
      },
      () => {
        console.log('定位失败，使用默认位置')
      }
    )
  }
})
</script>

<style scoped>
.home-container {
  display: flex;
  height: 100%;
  overflow: hidden;
  position: relative;
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

.search-panel {
  width: 350px;
  background: var(--bg-primary);
  overflow: hidden;
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  transition: transform 0.3s ease;
}

.search-panel-header {
  display: none;
  padding: 15px 20px;
  border-bottom: 1px solid var(--border-color);
  align-items: center;
  justify-content: space-between;
  background: var(--bg-secondary);
}

.search-panel-header h2 {
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
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: var(--transition-base);
}

.close-sidebar-btn:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.search-panel-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  /* 移动端滚动优化 */
  -webkit-overflow-scrolling: touch;
}

.search-panel h2 {
  margin: 0 0 20px 0;
  font-size: 20px;
  color: var(--text-primary);
}

.search-form {
  margin-bottom: 20px;
}

.form-row {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group-inline {
  margin-bottom: 0;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  color: var(--text-secondary);
  font-size: 14px;
}

.input {
  width: 100%;
  padding: 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  font-size: 14px;
  box-sizing: border-box;
  background: var(--bg-secondary);
  color: var(--text-primary);
  transition: var(--transition-base);
  /* 移动端优化 */
  -webkit-appearance: none;
  appearance: none;
}

.input:focus {
  outline: none;
  border-color: var(--border-color-focus);
  box-shadow: 0 0 0 2px rgba(139, 111, 71, 0.1);
}

.button-group {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

/* 桌面端保持垂直布局，移动端在媒体查询中覆盖 */
.button-group-row {
  flex-direction: column;
  gap: 10px;
}

.btn {
  padding: 12px 20px;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
  width: 100%;
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

.results-list {
  margin-top: 20px;
}

.results-list h3 {
  margin: 0 0 15px 0;
  font-size: 16px;
  color: var(--text-primary);
}

.result-item {
  background: var(--bg-secondary);
  padding: 15px;
  margin-bottom: 10px;
  border-radius: 4px;
  cursor: pointer;
  transition: var(--transition-base);
  border: 1px solid var(--border-color);
  /* 移动端触摸优化 */
  -webkit-tap-highlight-color: transparent;
}

.result-item:hover {
  box-shadow: var(--shadow-md);
  border-color: var(--border-color-hover);
  transform: translateY(-2px);
}

.result-item:active {
  transform: scale(0.98);
}

.result-content h4 {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: var(--text-primary);
}

.result-content p {
  margin: 5px 0;
  color: var(--text-secondary);
  font-size: 14px;
}

.distance {
  color: var(--color-primary) !important;
  font-weight: 500;
}

.result-image {
  width: 100%;
  max-width: 200px;
  height: auto;
  margin-top: 10px;
  border-radius: 4px;
}

.map-panel {
  flex: 1;
  position: relative;
  min-width: 0; /* 防止flex子元素溢出 */
}

/* 响应式设计 */
@media (max-width: 768px) {
  .home-container {
    flex-direction: column;
    height: 100vh;
    overflow: hidden;
  }

  .mobile-sidebar-toggle {
    display: none; /* 移动端不再需要侧边栏切换按钮 */
  }

  .sidebar-overlay {
    display: none; /* 移动端不再需要遮罩层 */
  }

  .search-panel {
    width: 100%;
    height: auto;
    max-height: 25vh;
    border-right: none;
    border-bottom: 1px solid var(--border-color);
    position: relative;
    transform: none;
    box-shadow: none;
    z-index: auto;
  }

  .search-panel.open {
    transform: none;
  }

  .search-panel-header {
    display: flex;
    padding: 10px 15px;
  }

  .search-panel-header h2 {
    font-size: 16px;
  }

  .search-panel-content {
    padding: 10px 15px;
    max-height: calc(25vh - 50px);
    overflow-y: auto;
  }

  .search-form {
    margin-bottom: 10px;
  }

  .form-row {
    flex-direction: row;
    gap: 10px;
    margin-bottom: 10px;
  }

  .form-group-inline {
    flex: 1;
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .form-group-inline label {
    font-size: 12px;
    margin-bottom: 0;
    white-space: nowrap;
    flex-shrink: 0;
    min-width: fit-content;
  }

  .form-group-inline .input {
    flex: 1;
    padding: 8px;
    font-size: 14px;
    min-width: 0;
  }

  .button-group-row {
    flex-direction: row;
    gap: 8px;
  }

  .button-group-row .btn {
    flex: 1;
    padding: 10px 8px;
    font-size: 13px;
    min-height: 40px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .map-panel {
    flex: 1;
    width: 100%;
    min-height: 75vh;
    height: 75vh;
  }
}

@media (max-width: 480px) {
  .search-panel {
    width: 100%;
    max-height: 25vh;
  }

  .search-panel-content {
    padding: 8px 12px;
    max-height: calc(25vh - 50px);
  }

  .search-panel-header {
    padding: 8px 12px;
  }

  .form-row {
    gap: 8px;
    margin-bottom: 8px;
  }

  .form-group-inline {
    gap: 4px;
  }

  .form-group-inline label {
    font-size: 11px;
  }

  .form-group-inline .input {
    padding: 6px 8px;
    font-size: 13px;
  }

  .button-group-row {
    flex-direction: row;
    gap: 6px;
  }

  .button-group-row .btn {
    flex: 1;
    padding: 8px 6px;
    font-size: 12px;
    min-height: 38px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .map-panel {
    min-height: 75vh;
    height: 75vh;
  }

  .result-item {
    padding: 12px;
  }
}

/* 手机端搜索结果按钮 */
.mobile-results-button {
  display: none;
  padding: 10px 15px;
  background: var(--bg-secondary);
  border-top: 1px solid var(--border-color);
}

/* 手机端搜索结果弹框 */
.mobile-results-modal {
  display: none;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 2000;
  pointer-events: none;
}

.mobile-results-modal.open {
  pointer-events: auto;
}

.modal-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.mobile-results-modal.open .modal-overlay {
  opacity: 1;
}

.modal-content {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: var(--bg-primary);
  border-radius: 20px 20px 0 0;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  transform: translateY(100%);
  transition: transform 0.3s ease;
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.2);
}

.mobile-results-modal.open .modal-content {
  transform: translateY(0);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px;
  border-bottom: 1px solid var(--border-color);
  flex-shrink: 0;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  color: var(--text-primary);
  font-weight: 600;
}

.modal-close-btn {
  background: none;
  border: none;
  font-size: 32px;
  color: var(--text-secondary);
  cursor: pointer;
  padding: 0;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: var(--transition-base);
  min-width: 44px;
  min-height: 44px;
}

.modal-close-btn:active {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 15px;
  -webkit-overflow-scrolling: touch;
}

.modal-body .result-item {
  margin-bottom: 12px;
}

.modal-body .result-item:last-child {
  margin-bottom: 0;
}

/* 手机端显示相关样式 */
@media (max-width: 768px) {
  .mobile-results-button {
    display: block;
  }
  
  .mobile-results-modal {
    display: block;
  }
  
  .results-list {
    display: none;
  }
}

/* 导航信息面板 */
.navigation-panel {
  position: absolute;
  bottom: 20px;
  left: 20px;
  right: 20px;
  max-width: 400px;
  background: var(--bg-primary);
  border-radius: 8px;
  box-shadow: var(--shadow-lg);
  z-index: 1000;
  border: 1px solid var(--border-color);
  overflow: hidden;
}

.navigation-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 15px 20px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
}

.navigation-header h4 {
  margin: 0;
  font-size: 16px;
  color: var(--text-primary);
  font-weight: 600;
}

.close-nav-btn {
  background: none;
  border: none;
  font-size: 24px;
  color: var(--text-secondary);
  cursor: pointer;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: var(--transition-base);
  line-height: 1;
}

.close-nav-btn:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}

.navigation-content {
  padding: 15px 20px;
}

.nav-info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.nav-info-item:last-of-type {
  margin-bottom: 15px;
}

.nav-label {
  color: var(--text-secondary);
  font-size: 14px;
}

.nav-value {
  color: var(--text-primary);
  font-size: 14px;
  font-weight: 500;
}

.nav-actions {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid var(--border-color);
}

.btn-sm {
  padding: 10px 16px;
  font-size: 13px;
}

/* 移动端导航面板样式 */
@media (max-width: 768px) {
  .navigation-panel {
    left: 10px;
    right: 10px;
    bottom: 10px;
    max-width: none;
  }

  .navigation-header {
    padding: 12px 15px;
  }

  .navigation-header h4 {
    font-size: 15px;
  }

  .navigation-content {
    padding: 12px 15px;
  }

  .nav-info-item {
    margin-bottom: 10px;
  }

  .nav-label,
  .nav-value {
    font-size: 13px;
  }

  .btn-sm {
    padding: 12px 16px;
    font-size: 14px;
    min-height: 44px;
  }
}

</style>

