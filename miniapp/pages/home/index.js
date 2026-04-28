const { getNearbyTrashCans, toggleLike, toggleDislike, getAllTags } = require("../../utils/api/trashcan")
const { resolveImageURL } = require("../../utils/request")
const { isAuthenticated } = require("../../utils/auth")
const { getAddressByLocation } = require("../../utils/geocoder")

Page({
  data: {
    latitude: 39.90923,
    longitude: 116.397428,
    scale: 15,
    radius: 2,
    limit: 20,
    loading: false,
    locating: false,
    trashCans: [],
    markers: [],
    allTags: [],
    selectedTagIndex: 0,
    selectedItem: null,
    showTip: false,
    currentAddress: ""
  },

  onLoad() {
    this.loadTags()
    this.locateMe(false).catch(() => {})
  },

  onShow() {
    if (this.data.trashCans.length > 0) {
      // already loaded, no need to reload
    }
  },

  loadTags() {
    getAllTags()
      .then((res) => {
        const tags = res.data || []
        this.setData({ allTags: tags })
      })
      .catch(() => {})
  },

  locateMe(showToast = true) {
    this.setData({ locating: true, showTip: false })

    return new Promise((resolve, reject) => {
      wx.getLocation({
        type: "gcj02",
        success: (res) => {
          const lat = res.latitude
          const lng = res.longitude

          this.setData({
            latitude: lat,
            longitude: lng
          })

          // reverse geocode
          getAddressByLocation(lng, lat)
            .then((address) => {
              if (address) {
                this.setData({ currentAddress: address })
              }
            })
            .catch(() => {})

          if (showToast) {
            wx.showToast({ title: "定位成功", icon: "success" })
          }

          this.searchNearby().finally(resolve)
        },
        fail: () => {
          wx.showToast({ title: "定位失败，请检查权限", icon: "none" })
          this.setData({ showTip: true })
          reject(new Error("location failed"))
        },
        complete: () => {
          this.setData({ locating: false })
        }
      })
    })
  },

  searchNearby() {
    const { latitude, longitude, radius, limit, selectedTagIndex, allTags } = this.data
    if (!latitude || !longitude) return

    this.setData({ loading: true })

    let tagId = ""
    if (selectedTagIndex > 0) {
      tagId = allTags[selectedTagIndex - 1]?.id || ""
    }

    return getNearbyTrashCans(latitude, longitude, radius, limit, tagId)
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
          ].filter((url) => url),
          distance_text: this.formatDistance(item.distance),
          walk_time: this.formatWalkTime(item.distance),
          user_action: Number(item.user_action || 0),
          like_count: Number(item.like_count || 0),
          dislike_count: Number(item.dislike_count || 0)
        }))

        this.setData({
          trashCans: list,
          markers: this.buildMarkers(list),
          showTip: list.length === 0
        })
      })
      .catch(() => {
        this.setData({ showTip: true })
      })
      .finally(() => {
        this.setData({ loading: false })
      })
  },

  buildMarkers(list) {
    return list
      .filter((item) => item.latitude && item.longitude)
      .map((item) => ({
        id: Number(item.id),
        latitude: Number(item.latitude),
        longitude: Number(item.longitude),
        width: 32,
        height: 32,
        iconPath: "/assets/marker.png",
        callout: {
          content: item.address || "垃圾桶",
          color: "#223126",
          fontSize: 12,
          borderRadius: 6,
          padding: 6,
          bgColor: "#ffffff",
          display: "BYCLICK"
        }
      }))
  },

  formatDistance(distance) {
    if (distance === undefined || distance === null) return "—"
    if (distance < 1) return `${Math.round(distance * 1000)}m`
    return `${distance.toFixed(1)}km`
  },

  formatWalkTime(distance) {
    if (distance === undefined || distance === null) return "—"
    const km = distance < 1 ? distance * 1000 / 1000 : distance
    const minutes = Math.round(km * 12)
    if (minutes < 1) return "1分钟"
    return `${minutes}分钟`
  },

  handleMarkerTap(e) {
    const markerId = Number(e.detail.markerId)
    const item = this.data.trashCans.find((it) => Number(it.id) === markerId)
    if (!item) return

    this.setData({
      selectedItem: item,
      latitude: Number(item.latitude),
      longitude: Number(item.longitude)
    })
  },

  openLocationPicker() {
    wx.chooseLocation({
      success: (res) => {
        if (!res.latitude || !res.longitude) return
        const lat = res.latitude
        const lng = res.longitude
        this.setData({
          latitude: lat,
          longitude: lng,
          currentAddress: res.address || res.name || ""
        })
        this.searchNearby()
      }
    })
  },

  handleRegionChange() {
    // optionally re-search on pan/zoom
  },

  handleLike(e) {
    if (!isAuthenticated()) {
      this.goLogin()
      return
    }
    const id = Number(e.currentTarget.dataset.id)
    const item = this.data.selectedItem
    if (!item || Number(item.id) !== id) return
    this.toggleLikeByItem(item)
  },

  handleDislike(e) {
    if (!isAuthenticated()) {
      this.goLogin()
      return
    }
    const id = Number(e.currentTarget.dataset.id)
    const item = this.data.selectedItem
    if (!item || Number(item.id) !== id) return
    this.toggleDislikeByItem(item)
  },

  toggleLikeByItem(item) {
    const action = item.user_action === 1 ? 0 : 1
    toggleLike(item.id)
      .then((res) => {
        this.patchSelectedItem(res.data || {})
      })
      .catch((err) => {
        wx.showToast({ title: err.message || "操作失败", icon: "none" })
      })
  },

  toggleDislikeByItem(item) {
    const action = item.user_action === -1 ? 0 : -1
    toggleDislike(item.id)
      .then((res) => {
        this.patchSelectedItem(res.data || {})
      })
      .catch((err) => {
        wx.showToast({ title: err.message || "操作失败", icon: "none" })
      })
  },

  patchSelectedItem(payload) {
    const item = this.data.selectedItem
    if (!item) return

    const updated = {
      ...item,
      like_count: Number(payload.like_count ?? item.like_count),
      dislike_count: Number(payload.dislike_count ?? item.dislike_count),
      user_action: Number(payload.user_action ?? item.user_action)
    }

    // also update in trashCans list
    const trashCans = this.data.trashCans.map((t) =>
      Number(t.id) === Number(item.id) ? updated : t
    )

    this.setData({
      selectedItem: updated,
      trashCans
    })
  },

  previewImage(e) {
    const urls = e.currentTarget.dataset.urls || []
    const current = e.currentTarget.dataset.current || ""
    if (!urls.length) return
    wx.previewImage({ urls, current })
  },

  navigateToItem(e) {
    const id = e.currentTarget.dataset.id
    const item = this.data.trashCans.find((t) => Number(t.id) === Number(id))
    if (!item) return
    wx.openLocation({
      latitude: Number(item.latitude),
      longitude: Number(item.longitude),
      scale: 18,
      name: "垃圾桶位置",
      address: item.address || ""
    })
  },

  goLogin() {
    wx.navigateTo({ url: "/pages/login/index" })
  }
})
