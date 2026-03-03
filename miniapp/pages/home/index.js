const { getNearbyTrashCans, toggleLike, toggleDislike } = require("../../utils/api/trashcan")
const { resolveImageURL } = require("../../utils/request")
const { isAuthenticated } = require("../../utils/auth")

Page({
  data: {
    latitude: 39.90923,
    longitude: 116.397428,
    scale: 15,
    radius: 1,
    limit: 10,
    loading: false,
    locating: false,
    authed: false,
    userLocation: null,
    trashCans: [],
    markers: []
  },

  onLoad() {
    this.setData({ authed: isAuthenticated() })
    this.locateMe(false).catch(() => {})
  },

  onShow() {
    this.setData({ authed: isAuthenticated() })
  },

  onPullDownRefresh() {
    if (!this.data.userLocation) {
      this.locateMe(false)
        .catch(() => {})
        .finally(() => wx.stopPullDownRefresh())
      return
    }

    this.searchNearby().finally(() => wx.stopPullDownRefresh())
  },

  handleRadiusInput(e) {
    const value = Number(e.detail.value || 1)
    this.setData({ radius: value > 0 ? value : 1 })
  },

  handleLimitInput(e) {
    const value = Number(e.detail.value || 10)
    this.setData({ limit: value > 0 ? value : 10 })
  },

  locateMe(showToast = true) {
    this.setData({ locating: true })

    return new Promise((resolve, reject) => {
      wx.getLocation({
        type: "gcj02",
        success: (res) => {
          const userLocation = {
            latitude: res.latitude,
            longitude: res.longitude
          }

          this.setData({
            latitude: res.latitude,
            longitude: res.longitude,
            userLocation
          })

          if (showToast) {
            wx.showToast({
              title: "定位成功",
              icon: "success"
            })
          }

          this.searchNearby().finally(resolve)
        },
        fail: (err) => {
          wx.showToast({
            title: "定位失败，请检查权限",
            icon: "none"
          })
          reject(err)
        },
        complete: () => {
          this.setData({ locating: false })
        }
      })
    })
  },

  searchNearby() {
    const { userLocation, radius, limit } = this.data

    if (!userLocation) {
      wx.showToast({
        title: "请先定位",
        icon: "none"
      })
      return Promise.resolve()
    }

    this.setData({ loading: true })

    return getNearbyTrashCans(
      userLocation.latitude,
      userLocation.longitude,
      radius,
      limit
    )
      .then((res) => {
        const list = (res.data || []).map((item) => ({
          ...item,
          image_url: resolveImageURL(item.image_url),
          image_url_2: resolveImageURL(item.image_url_2),
          image_url_3: resolveImageURL(item.image_url_3),
          image_urls: [
            resolveImageURL(item.image_url),
            resolveImageURL(item.image_url_2),
            resolveImageURL(item.image_url_3)
          ].filter(url => url),
          distance_text: this.formatDistance(item.distance),
          user_action: Number(item.user_action || 0),
          like_count: Number(item.like_count || 0),
          dislike_count: Number(item.dislike_count || 0)
        }))

        this.setData({
          trashCans: list,
          markers: this.buildMarkers(list)
        })
      })
      .catch((err) => {
        wx.showToast({
          title: err.message || "搜索失败",
          icon: "none"
        })
      })
      .finally(() => {
        this.setData({ loading: false })
      })
  },

  buildMarkers(list) {
    return list
      .filter((item) => item.latitude && item.longitude)
      .map((item) => {
        const uploaderText = item.uploader_name ? `👤 ${item.uploader_name}` : ""
        const addressText = item.address || "垃圾桶"
        const content = uploaderText 
          ? `${uploaderText}\n${addressText}\n${this.formatDistance(item.distance)}`
          : `${addressText}\n${this.formatDistance(item.distance)}`
        
        return {
          id: Number(item.id),
          latitude: Number(item.latitude),
          longitude: Number(item.longitude),
          width: 28,
          height: 28,
          callout: {
            content,
            color: "#223126",
            fontSize: 12,
            borderRadius: 6,
            padding: 6,
            bgColor: "#ffffff",
            display: "BYCLICK"
          }
        }
      })
  },

  formatDistance(distance) {
    if (distance === undefined || distance === null) {
      return "距离未知"
    }

    return `距离 ${Number(distance).toFixed(2)} km`
  },

  handleMarkerTap(e) {
    const markerId = Number(e.detail.markerId)
    const item = this.data.trashCans.find((it) => Number(it.id) === markerId)
    if (!item) {
      return
    }

    this.setData({
      latitude: Number(item.latitude),
      longitude: Number(item.longitude)
    })

    this.openActionSheet(item)
  },

  openActionSheet(item) {
    const actions = ["导航到这里"]
    if (item.image_urls && item.image_urls.length > 0) {
      actions.push("查看图片")
    }

    actions.push(`点赞 (${item.like_count || 0})`)
    actions.push(`点踩 (${item.dislike_count || 0})`)

    wx.showActionSheet({
      itemList: actions,
      success: (res) => {
        const action = actions[res.tapIndex]

        if (action === "导航到这里") {
          this.navigateToItemByData(item)
          return
        }

        if (action === "查看图片") {
          wx.previewImage({
            urls: item.image_urls || [],
            current: item.image_urls ? item.image_urls[0] : ''
          })
          return
        }

        if (!this.ensureAuth()) {
          return
        }

        if (action.startsWith("点赞")) {
          this.handleLikeByData(item)
          return
        }

        if (action.startsWith("点踩")) {
          this.handleDislikeByData(item)
        }
      }
    })
  },

  handleLike(e) {
    const id = Number(e.currentTarget.dataset.id)
    const item = this.data.trashCans.find((it) => Number(it.id) === id)
    if (!item) {
      return
    }

    this.handleLikeByData(item)
  },

  handleDislike(e) {
    const id = Number(e.currentTarget.dataset.id)
    const item = this.data.trashCans.find((it) => Number(it.id) === id)
    if (!item) {
      return
    }

    this.handleDislikeByData(item)
  },

  handleLikeByData(item) {
    if (!this.ensureAuth()) {
      return
    }

    toggleLike(item.id)
      .then((res) => {
        this.patchActionResult(item.id, res.data || {})
      })
      .catch((err) => {
        wx.showToast({
          title: err.message || "操作失败",
          icon: "none"
        })
      })
  },

  handleDislikeByData(item) {
    if (!this.ensureAuth()) {
      return
    }

    toggleDislike(item.id)
      .then((res) => {
        this.patchActionResult(item.id, res.data || {})
      })
      .catch((err) => {
        wx.showToast({
          title: err.message || "操作失败",
          icon: "none"
        })
      })
  },

  patchActionResult(id, payload) {
    const list = this.data.trashCans.map((item) => {
      if (Number(item.id) !== Number(id)) {
        return item
      }

      return {
        ...item,
        distance_text: this.formatDistance(item.distance),
        like_count: Number(payload.like_count || 0),
        dislike_count: Number(payload.dislike_count || 0),
        user_action: Number(payload.user_action || 0)
      }
    })

    this.setData({
      trashCans: list,
      markers: this.buildMarkers(list)
    })
  },

  previewImage(e) {
    const urls = e.currentTarget.dataset.urls
    const current = e.currentTarget.dataset.current
    if (!urls || urls.length === 0) {
      return
    }

    wx.previewImage({
      urls: urls,
      current: current
    })
  },

  navigateToItem(e) {
    const id = Number(e.currentTarget.dataset.id)
    const item = this.data.trashCans.find((it) => Number(it.id) === id)
    if (!item) {
      return
    }

    this.navigateToItemByData(item)
  },

  navigateToItemByData(item) {
    wx.openLocation({
      latitude: Number(item.latitude),
      longitude: Number(item.longitude),
      scale: 18,
      name: "垃圾桶位置",
      address: item.address || ""
    })
  },

  ensureAuth() {
    if (isAuthenticated()) {
      return true
    }

    wx.showModal({
      title: "需要登录",
      content: "该操作需要先登录，是否去登录？",
      success: (res) => {
        if (res.confirm) {
          const redirect = encodeURIComponent("/pages/home/index")
          wx.navigateTo({
            url: `/pages/login/index?redirect=${redirect}`
          })
        }
      }
    })

    return false
  },

  goUpload() {
    wx.navigateTo({
      url: "/pages/upload/index"
    })
  },

  goProfileOrLogin() {
    if (isAuthenticated()) {
      wx.navigateTo({
        url: "/pages/profile/index"
      })
      return
    }

    wx.navigateTo({
      url: "/pages/login/index?redirect=%2Fpages%2Fprofile%2Findex"
    })
  }
})
