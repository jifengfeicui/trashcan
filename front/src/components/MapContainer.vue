<template>
  <div id="map-container" ref="mapContainer"></div>
</template>

<script setup>
import {nextTick, onMounted, onUnmounted, ref, watch} from 'vue'

const props = defineProps({
  center: {
    type: Array,
    default: () => [121.4375, 31.1956] // 默认上海徐家汇
  },
  markers: {
    type: Array,
    default: () => []
  },
  zoom: {
    type: Number,
    default: 13
  }
})

const emit = defineEmits(['map-click', 'marker-click', 'map-ready', 'location-ready', 'info-window-action', 'navigation-ready'])

const mapContainer = ref(null)
let map = null
// 使用 Map 存储 marker 和 infoWindow，便于管理生命周期
const markerMap = new Map() // id -> { marker, infoWindow, trashCan }
let userLocationMarker = null // 用户位置标记
let userLocation = null // 用户位置坐标 { lng, lat }
let navigationInstance = null // 导航实例
let navigationRoute = null // 导航路线
let globalClickHandler = null // 全局点击事件处理器

// 等待 AMap 加载完成
const waitForAMap = () => {
  return new Promise((resolve, reject) => {
    if (window.AMap) {
      resolve(window.AMap)
      return
    }

    // 如果已经在加载中，等待
    let checkCount = 0
    const checkInterval = setInterval(() => {
      checkCount++
      if (window.AMap) {
        clearInterval(checkInterval)
        resolve(window.AMap)
      } else if (checkCount > 50) { // 最多等待 5 秒
        clearInterval(checkInterval)
        reject(new Error('高德地图API加载超时，请检查网络连接和API Key配置'))
      }
    }, 100)
  })
}

// 初始化地图
const initMap = async () => {
  try {
    // 确保容器已经渲染
    if (!mapContainer.value) {
      console.error('地图容器未找到')
      return
    }

    // 等待 AMap 加载完成（从 index.html 引入）
    const AMap = await waitForAMap()

    // 等待下一帧确保容器有尺寸
    await new Promise(resolve => setTimeout(resolve, 100))

    // 检查容器尺寸
    const containerRect = mapContainer.value.getBoundingClientRect()
    console.log('地图容器尺寸:', containerRect.width, 'x', containerRect.height)

    if (containerRect.width === 0 || containerRect.height === 0) {
      console.warn('地图容器尺寸为0，等待容器渲染...')
      await new Promise(resolve => setTimeout(resolve, 500))
    }

    map = new AMap.Map(mapContainer.value, {
      zoom: props.zoom,
      center: props.center,
      viewMode: '2D',
      mapStyle: 'amap://styles/normal'
    })

    map.on('error', (e) => {
      console.error('地图加载错误:', e)
    })

    // 地图加载完成事件 - 只代表地图初始化完成
    map.on('complete', () => {
      console.log('地图加载完成')
      emit('map-ready')
    })

    // 地图点击事件
    map.on('click', (e) => {
      emit('map-click', {
        lng: e.lnglat.getLng(),
        lat: e.lnglat.getLat()
      })
    })

    // 获取用户当前位置
    map.plugin('AMap.Geolocation', () => {
      const geolocation = new AMap.Geolocation({
        enableHighAccuracy: true,
        timeout: 10000,
        buttonOffset: new AMap.Pixel(10, 20),
        zoomToAccuracy: true,
        buttonPosition: 'RB',
        showMarker: false // 禁用默认的蓝色三角标记，只使用自定义的蓝色圆点
      })

      map.addControl(geolocation)

      geolocation.getCurrentPosition((status, result) => {
        if (status === 'complete') {
          const {lng, lat} = result.position
          map.setCenter([lng, lat])
          // 使用 locateUser 方法显示位置标记
          locateUser(lng, lat)
        } else {
          console.error('定位失败:', result)
        }
      })
    })
  } catch (error) {
    console.error('地图加载失败:', error)
    const errorMsg = error.message || String(error)
    if (errorMsg.includes('USERKEY_PLAT_NOMATCH') || errorMsg.includes('平台')) {
      console.error('❌ API Key平台类型不匹配！')
      console.error('💡 解决方案：')
      console.error('1. 登录高德开放平台：https://console.amap.com/')
      console.error('2. 进入"应用管理" -> 找到你的Key')
      console.error('3. 确保"服务平台"设置为"Web端(JS API)"')
      console.error('4. 如果设置了白名单，请添加当前域名或IP（localhost、127.0.0.1等）')
      alert('地图加载失败：API Key配置错误\n\n请检查：\n1. API Key是否设置为"Web端(JS API)"\n2. 是否设置了域名/IP白名单限制\n\n详细错误信息请查看控制台')
    } else {
      alert('地图加载失败：' + errorMsg + '\n\n请查看控制台获取详细信息')
    }
  }
}

// 清除所有标记
const clearMarkers = () => {
  markerMap.forEach(({marker, infoWindow}) => {
    infoWindow.close()
    marker.setMap(null) // 正确清理 marker
  })
  markerMap.clear()
}

// 定位到用户位置并显示标记
const locateUser = (lng, lat) => {
  if (!map || !window.AMap) {
    console.error('地图未初始化')
    return
  }

  // 清除之前的用户位置标记
  if (userLocationMarker) {
    userLocationMarker.setMap(null) // 正确清理
    userLocationMarker = null
  }

  // 创建用户位置标记（蓝色圆点）
  userLocationMarker = new window.AMap.Marker({
    position: [lng, lat],
    title: '我的位置',
    icon: new window.AMap.Icon({
      size: new window.AMap.Size(40, 40),
      image: 'https://webapi.amap.com/theme/v1.3/markers/n/mark_b.png',
      imageOffset: new window.AMap.Pixel(0, 0),
      imageSize: new window.AMap.Size(40, 40)
    }),
    zIndex: 1000 // 确保在最上层
  })

  map.add(userLocationMarker)

  // 保存用户位置
  userLocation = { lng, lat }

  // 设置地图中心并调整缩放级别
  map.setCenter([lng, lat])
  map.setZoom(16)

  // 输出详细位置信息到控制台（方便伪造数据测试）
  const locationInfo = {
    longitude: lng,
    latitude: lat,
    lng: lng,
    lat: lat,
    coords: [lng, lat],
    // 格式化的位置信息，方便复制
    formatted: {
      'WGS84坐标': `${lat}, ${lng}`,
      'GCJ02坐标（高德）': `${lat}, ${lng}`,
      '数组格式': `[${lng}, ${lat}]`,
      '对象格式': `{ lng: ${lng}, lat: ${lat} }`,
      'JSON格式': JSON.stringify({lng, lat})
    }
  }

  console.log('========== 位置信息 ==========')
  console.log('经度 (longitude/lng):', lng)
  console.log('纬度 (latitude/lat):', lat)
  console.log('坐标数组:', [lng, lat])
  console.log('坐标对象:', {lng, lat})
  console.log('格式化信息:', locationInfo.formatted)
  console.log('完整位置对象:', locationInfo)
  console.log('============================')

  // 触发位置就绪事件（与地图就绪事件分离）
  emit('location-ready', {lng, lat})

  return locationInfo
}

// 获取图片完整URL
const getImageUrl = (imagePath) => {
  if (!imagePath) return ''
  if (imagePath.startsWith('http://') || imagePath.startsWith('https://')) {
    return imagePath
  }
  return imagePath.startsWith('/') ? imagePath : `/${imagePath}`
}

// 处理 InfoWindow 中的操作（通过事件而非全局函数）
const handleInfoWindowAction = (action, data) => {
  emit('info-window-action', {action, data})
}

// 添加标记
const addMarker = (trashCan) => {
  if (!map || !window.AMap) return

  const {latitude, longitude, address, description, image_url, distance, id, like_count = 0, dislike_count = 0, user_action = 0} = trashCan

  // 如果已存在，先移除
  if (markerMap.has(id)) {
    removeMarker(id)
  }

  // 创建标记
  const marker = new window.AMap.Marker({
    position: [longitude, latitude],
    title: address || '垃圾桶',
    icon: new window.AMap.Icon({
      size: new window.AMap.Size(32, 32),
      image: 'https://webapi.amap.com/theme/v1.3/markers/n/mark_r.png',
      imageOffset: new window.AMap.Pixel(0, 0),
      imageSize: new window.AMap.Size(32, 32)
    })
  })

  // 处理图片URL（转义防止 XSS）
  const imageUrl = image_url ? getImageUrl(image_url) : ''
  const safeAddress = (address || '').replace(/</g, '&lt;').replace(/>/g, '&gt;')
  const safeDescription = (description || '').replace(/</g, '&lt;').replace(/>/g, '&gt;')

  // 创建信息窗口内容 - 使用 data-action 属性而非 onclick
  // 检测是否为移动设备
  const isMobile = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent) || window.innerWidth <= 768
  
  const actionId = `action_${id}_${Date.now()}`
  const infoWindowId = `infowindow_${id}`
  const buttonPadding = isMobile ? '12px 16px' : '6px 12px'
  const buttonFontSize = isMobile ? '16px' : '14px'
  const minWidth = isMobile ? '280px' : '200px'
  const maxImageWidth = isMobile ? '100%' : '300px'
  
  let content = `
    <div id="${infoWindowId}" style="padding: ${isMobile ? '15px' : '10px'}; min-width: ${minWidth}; max-width: ${isMobile ? '90vw' : '350px'};">
      <h3 style="margin: 0 0 ${isMobile ? '12px' : '10px'} 0; font-size: ${isMobile ? '18px' : '16px'}; color: #2C2416; font-weight: 600;">垃圾桶位置</h3>
      ${safeAddress ? `<p style="margin: ${isMobile ? '8px' : '5px'} 0; color: #5C4E3A; font-size: ${isMobile ? '15px' : '14px'}; line-height: 1.5;">${safeAddress}</p>` : ''}
      ${safeDescription ? `<p style="margin: ${isMobile ? '8px' : '5px'} 0; color: #5C4E3A; font-size: ${isMobile ? '15px' : '14px'}; line-height: 1.5;">${safeDescription}</p>` : ''}
      ${distance !== undefined ? `<p style="margin: ${isMobile ? '8px' : '5px'} 0; color: #8B6F47; font-weight: 500; font-size: ${isMobile ? '15px' : '14px'};">距离: ${distance.toFixed(2)} 公里</p>` : ''}
      ${imageUrl ? `<img src="${imageUrl}" style="width: 100%; max-width: ${maxImageWidth}; margin-top: ${isMobile ? '12px' : '10px'}; border-radius: 4px; cursor: pointer; display: block;" alt="垃圾桶图片" data-action="open-image" data-image-url="${imageUrl.replace(/"/g, '&quot;')}" />` : ''}
      <div style="margin-top: ${isMobile ? '15px' : '10px'}; display: flex; gap: ${isMobile ? '10px' : '8px'}; flex-wrap: wrap;">
        <button data-action="like" data-id="${id}" 
                style="flex: 1; background: ${user_action === 1 ? '#8B6F47' : '#f5f5f5'}; color: ${user_action === 1 ? 'white' : '#333'}; border: 1px solid ${user_action === 1 ? '#8B6F47' : '#ddd'}; padding: ${buttonPadding}; border-radius: 4px; cursor: pointer; transition: all 0.3s; font-size: ${buttonFontSize}; min-height: ${isMobile ? '44px' : 'auto'}; font-weight: 500; -webkit-tap-highlight-color: transparent;">
          👍 ${like_count || 0}
        </button>
        <button data-action="dislike" data-id="${id}" 
                style="flex: 1; background: ${user_action === -1 ? '#8B6F47' : '#f5f5f5'}; color: ${user_action === -1 ? 'white' : '#333'}; border: 1px solid ${user_action === -1 ? '#8B6F47' : '#ddd'}; padding: ${buttonPadding}; border-radius: 4px; cursor: pointer; transition: all 0.3s; font-size: ${buttonFontSize}; min-height: ${isMobile ? '44px' : 'auto'}; font-weight: 500; -webkit-tap-highlight-color: transparent;">
          👎 ${dislike_count || 0}
        </button>
      </div>
      <div style="margin-top: ${isMobile ? '12px' : '10px'};">
        <button data-action="navigate" data-lng="${longitude}" data-lat="${latitude}" 
                style="background: #8B6F47; color: white; border: none; padding: ${buttonPadding}; border-radius: 4px; cursor: pointer; transition: all 0.3s; font-size: ${buttonFontSize}; width: 100%; min-height: ${isMobile ? '44px' : 'auto'}; font-weight: 500; -webkit-tap-highlight-color: transparent;">
          🚶 步行路线
        </button>
      </div>
    </div>
  `

  const infoWindow = new window.AMap.InfoWindow({
    content: content,
    offset: new window.AMap.Pixel(0, -30)
  })

  // 标记点击事件
  marker.on('click', () => {
    // 关闭其他信息窗口
    markerMap.forEach(({infoWindow: iw}) => {
      if (iw !== infoWindow) {
        iw.close()
      }
    })

    infoWindow.open(map, marker.getPosition())

    // 使用全局事件委托处理 InfoWindow 内的操作（更可靠的方法）
    // 延迟绑定以确保 DOM 已渲染
    setTimeout(() => {
      // 查找 InfoWindow 内的按钮元素
      const findButton = () => {
        // 方法1: 通过 ID 查找容器
        let containerEl = document.getElementById(infoWindowId)
        
        // 方法2: 通过 data-action 和坐标查找按钮
        if (!containerEl) {
          const buttons = document.querySelectorAll(`[data-action="navigate"][data-lng="${longitude}"][data-lat="${latitude}"]`)
          if (buttons.length > 0) {
            containerEl = buttons[0].closest('div') || buttons[0].parentElement
          }
        }
        
        // 方法3: 尝试通过 getContent() 获取
        if (!containerEl) {
          const infoWindowContent = infoWindow.getContent()
          if (infoWindowContent) {
            if (infoWindowContent.nodeType === 1) {
              containerEl = infoWindowContent
            } else if (infoWindowContent.querySelector) {
              containerEl = infoWindowContent.querySelector(`#${infoWindowId}`) || infoWindowContent
            }
          }
        }
        
        return containerEl
      }

      const containerEl = findButton()

      if (!containerEl || typeof containerEl.addEventListener !== 'function') {
        console.warn('InfoWindow 内容元素无效，无法绑定事件，将使用全局事件委托', {
          containerEl,
          infoWindowId,
          longitude,
          latitude
        })
        
        // 如果找不到容器，使用全局事件委托（作为后备方案）
        // 注意：这已经在组件级别处理，不需要在这里重复
        return
      }

      const handleClick = (e) => {
        // 使用 closest 查找带有 data-action 的元素（处理按钮内文本节点的情况）
        const target = e.target.closest ? e.target.closest('[data-action]') : (() => {
          // Fallback: 向上查找带有 data-action 的元素
          let el = e.target
          while (el && el !== containerEl) {
            if (el.getAttribute && el.getAttribute('data-action')) {
              return el
            }
            el = el.parentElement
          }
          return null
        })()

        if (!target || !target.getAttribute) return

        const action = target.getAttribute('data-action')

        if (action === 'navigate') {
          e.preventDefault()
          e.stopPropagation()
          const lng = parseFloat(target.getAttribute('data-lng'))
          const lat = parseFloat(target.getAttribute('data-lat'))
          console.log('步行路线按钮被点击，目标位置:', {lng, lat})
          handleInfoWindowAction('navigate', {lng, lat})
        } else if (action === 'open-image') {
          e.preventDefault()
          e.stopPropagation()
          const imageUrl = target.getAttribute('data-image-url')
          handleInfoWindowAction('open-image', {imageUrl})
        } else if (action === 'like' || action === 'dislike') {
          e.preventDefault()
          e.stopPropagation()
          const trashCanId = target.getAttribute('data-id')
          handleInfoWindowAction(action, {id: trashCanId})
        }
      }

      // 移除旧的事件监听器（如果有）
      const oldHandler = containerEl._clickHandler
      if (oldHandler) {
        containerEl.removeEventListener('click', oldHandler)
      }

      // 添加新的事件监听器（使用捕获阶段确保事件被捕获）
      containerEl.addEventListener('click', handleClick, true)
      containerEl._clickHandler = handleClick
      
      console.log('InfoWindow 事件监听器已绑定', containerEl)
    }, 200)

    emit('marker-click', trashCan)
  })

  map.add(marker)

  // 存储到 Map 中
  markerMap.set(id, {marker, infoWindow, trashCan})
}

// 移除单个标记
const removeMarker = (id) => {
  const item = markerMap.get(id)
  if (item) {
    item.infoWindow.close()
    item.marker.setMap(null) // 正确清理
    markerMap.delete(id)
  }
}

// 监听markers变化 - 使用 diff 机制避免性能问题
watch(() => props.markers, (newMarkers, oldMarkers) => {
  if (!map || !newMarkers) return

  const newIds = new Set(newMarkers.map(m => m.id))
  const oldIds = new Set((oldMarkers || []).map(m => m.id))

  // 移除不存在的标记
  oldIds.forEach(id => {
    if (!newIds.has(id)) {
      removeMarker(id)
    }
  })

  // 添加或更新标记
  const markerMapById = new Map(newMarkers.map(m => [m.id, m]))
  newMarkers.forEach(trashCan => {
    const existing = markerMap.get(trashCan.id)
    // 如果不存在或数据有变化，则添加/更新
    if (!existing || existing.trashCan !== trashCan) {
      addMarker(trashCan)
    }
  })
}, {immediate: false})

// 监听center变化
watch(() => props.center, (newCenter) => {
  if (map && newCenter && newCenter.length === 2) {
    map.setCenter(newCenter)
  }
})

onMounted(async () => {
  // 添加全局事件委托作为后备方案（处理 InfoWindow 内的按钮点击）
  globalClickHandler = (e) => {
    // 查找带有 data-action 的元素
    const target = e.target.closest ? e.target.closest('[data-action]') : (() => {
      let el = e.target
      while (el && el !== document.body) {
        if (el.getAttribute && el.getAttribute('data-action')) {
          return el
        }
        el = el.parentElement
      }
      return null
    })()

    if (!target || !target.getAttribute) return

    const action = target.getAttribute('data-action')

    // 检查是否点击了 InfoWindow 内的按钮（通过检查是否在高德地图的 InfoWindow 容器内）
    const isInInfoWindow = target.closest ? target.closest('.amap-info-content, .amap-info') : null

    if (isInInfoWindow) {
      if (action === 'navigate') {
        e.preventDefault()
        e.stopPropagation()
        const lng = parseFloat(target.getAttribute('data-lng'))
        const lat = parseFloat(target.getAttribute('data-lat'))
        console.log('步行路线按钮被点击（全局事件委托），目标位置:', {lng, lat})
        handleInfoWindowAction('navigate', {lng, lat})
      } else if (action === 'open-image') {
        e.preventDefault()
        e.stopPropagation()
        const imageUrl = target.getAttribute('data-image-url')
        handleInfoWindowAction('open-image', {imageUrl})
      } else if (action === 'like' || action === 'dislike') {
        e.preventDefault()
        e.stopPropagation()
        const trashCanId = target.getAttribute('data-id')
        handleInfoWindowAction(action, {id: trashCanId})
      }
    }
  }

  // 使用捕获阶段确保事件被捕获
  document.addEventListener('click', globalClickHandler, true)

  // 等待DOM渲染完成
  await nextTick()
  // 再等待一小段时间确保容器有尺寸
  setTimeout(() => {
    initMap()
  }, 200)
})

onUnmounted(() => {
  // 移除全局事件监听器
  if (globalClickHandler) {
    document.removeEventListener('click', globalClickHandler, true)
    globalClickHandler = null
  }

  // 清理导航路线
  clearNavigation()

  // 清理所有标记
  clearMarkers()

  // 清理用户位置标记
  if (userLocationMarker && map) {
    userLocationMarker.setMap(null)
    userLocationMarker = null
  }

  // 销毁地图
  if (map) {
    map.destroy()
    map = null
  }
})

// 清除导航路线
const clearNavigation = () => {
  if (navigationRoute) {
    map.remove(navigationRoute)
    navigationRoute = null
  }
  if (navigationInstance) {
    navigationInstance.clear()
    navigationInstance = null
  }
}

// 导航到指定位置（步行路线规划）
const navigateTo = async (destination, startPoint = null) => {
  if (!map || !window.AMap) {
    console.error('地图未初始化')
    return
  }

  // 清除之前的导航路线
  clearNavigation()

  // 获取起点（用户当前位置或指定起点）
  let start = startPoint
  if (!start && userLocation) {
    start = [userLocation.lng, userLocation.lat]
  } else if (!start && userLocationMarker) {
    const position = userLocationMarker.getPosition()
    start = [position.getLng(), position.getLat()]
  }

  if (!start) {
    console.error('无法获取起点位置')
    return Promise.reject(new Error('请先获取您的位置'))
  }

  const [startLng, startLat] = Array.isArray(start) ? start : [start.lng || start.longitude, start.lat || start.latitude]
  const [destLng, destLat] = Array.isArray(destination) ? destination : [destination.lng || destination.longitude, destination.lat || destination.latitude]

  return new Promise((resolve, reject) => {
    // 使用步行导航插件
    map.plugin('AMap.Walking', () => {
      try {
        if (!window.AMap.Walking) {
          reject(new Error('步行导航插件加载失败'))
          return
        }

        // 创建步行导航实例
        navigationInstance = new window.AMap.Walking({
          map: map,
          hideMarkers: false,
          autoFitView: true
        })

        // 搜索路径规划
        navigationInstance.search(
          new window.AMap.LngLat(startLng, startLat),
          new window.AMap.LngLat(destLng, destLat),
          (status, result) => {
            if (status === 'complete') {
              console.log('步行路线规划成功:', result)
              
              // 获取路线信息
              const route = result.routes && result.routes[0]
              if (route) {
                const distance = (route.distance / 1000).toFixed(2) // 转换为公里
                const duration = Math.round(route.time / 60) // 转换为分钟
                
                // 触发导航完成事件
                emit('navigation-ready', {
                  mode: 'walking',
                  distance: parseFloat(distance),
                  duration,
                  route: result.routes[0]
                })
              }

              resolve({
                mode: 'walking',
                distance: route ? (route.distance / 1000).toFixed(2) : 0,
                duration: route ? Math.round(route.time / 60) : 0,
                route: result.routes[0]
              })
            } else {
              console.error('步行路线规划失败:', result)
              reject(new Error(result.message || '路线规划失败'))
            }
          }
        )
      } catch (error) {
        console.error('步行导航功能初始化失败:', error)
        reject(error)
      }
    })
  })
}

// 暴露方法给父组件
defineExpose({
  addMarker,
  clearMarkers,
  removeMarker,
  locateUser,
  navigateTo,
  clearNavigation,
  getMap: () => map
})
</script>

<style scoped>
#map-container {
  width: 100%;
  height: 100%;
  min-height: 500px;
  position: relative;
  /* 移动端优化 */
  touch-action: pan-x pan-y;
}

/* 移动端优化 */
@media (max-width: 768px) {
  #map-container {
    min-height: 400px;
  }
}

@media (max-width: 480px) {
  #map-container {
    min-height: 300px;
  }
}
</style>

