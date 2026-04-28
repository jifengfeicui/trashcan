const { getNearbyTrashCans, toggleLike, toggleDislike, getAllTags } = require("../../utils/api/trashcan")
const { resolveImageURL } = require("../../utils/request")
const { isAuthenticated } = require("../../utils/auth")

Page({
  data: {
    latitude: "",
    longitude: "",
    radius: 2,
    limit: 20,
    loading: false,
    trashCans: [],
    allTags: [],
    tagOptions: ["全部"],
    selectedTagIndex: 0,
    selectedRadiusIndex: 0,
    radiusOptions: ["500m", "1km", "2km", "5km"],
    radiusMap: { "500m": 0.5, "1km": 1, "2km": 2, "5km": 5 }
  },

  onLoad() {
    this.loadTags()
    this.locateAndSearch()
  },

  onShow() {
    const app = getApp()
    if (app.globalData && app.globalData.needRefreshDiscover) {
      app.globalData.needRefreshDiscover = false
      this.locateAndSearch()
    }
  },

  loadTags() {
    getAllTags()
      .then((res) => {
        const tags = res.data || []
        const tagOptions = ["全部", ...tags.map((t) => t.name)]
        this.setData({ allTags: tags, tagOptions: tagOptions })
      })
      .catch(() => {})
  },

  locateAndSearch() {
    this.setData({ loading: true })
    wx.getLocation({
      type: "gcj02",
      success: (res) => {
        const lat = res.latitude
        const lng = res.longitude
        this.setData({ latitude: lat, longitude: lng })
        this.search()
      },
      fail: () => {
        wx.showToast({ title: "定位失败，请在设置中开启定位", icon: "none" })
        this.setData({ loading: false })
      }
    })
  },

  search() {
    const { latitude, longitude, radius, limit, selectedTagIndex, allTags } = this.data
    if (!latitude || !longitude) return

    this.setData({ loading: true })
    let tagId = ""
    if (selectedTagIndex > 0) {
      tagId = allTags[selectedTagIndex - 1]?.id || ""
    }

    getNearbyTrashCans(latitude, longitude, radius, limit, tagId)
      .then((res) => {
        const list = ((res.data || [])).map((item) => ({
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
          loading: false
        })
      })
      .catch(() => {
        this.setData({ loading: false })
      })
  },

  onRadiusChange(e) {
    const idx = Number(e.detail.value)
    const radiusOptions = this.data.radiusOptions
    const radius = this.data.radiusMap[radiusOptions[idx]]
    this.setData({ selectedRadiusIndex: idx, radius })
    this.search()
  },

  onTagChange(e) {
    const index = Number(e.detail.value)
    this.setData({ selectedTagIndex: index })
    this.search()
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

  handleLike(e) {
    if (!isAuthenticated()) {
      this.goLogin()
      return
    }
    const id = Number(e.currentTarget.dataset.id)
    const item = this.data.trashCans.find((t) => Number(t.id) === id)
    if (!item) return
    const action = item.user_action === 1 ? 0 : 1
    this.patchItem(id, { action })
    toggleLike(id)
      .then((res) => {
        this.patchItem(id, res.data || {})
      })
      .catch(() => {})
  },

  handleDislike(e) {
    if (!isAuthenticated()) {
      this.goLogin()
      return
    }
    const id = Number(e.currentTarget.dataset.id)
    const item = this.data.trashCans.find((t) => Number(t.id) === id)
    if (!item) return
    const action = item.user_action === -1 ? 0 : -1
    this.patchItem(id, { action })
    toggleDislike(id)
      .then((res) => {
        this.patchItem(id, res.data || {})
      })
      .catch(() => {})
  },

  patchItem(id, payload) {
    const trashCans = this.data.trashCans.map((t) =>
      Number(t.id) === Number(id)
        ? {
            ...t,
            like_count: Number(payload.like_count ?? t.like_count),
            dislike_count: Number(payload.dislike_count ?? t.dislike_count),
            user_action: Number(payload.user_action ?? t.user_action)
          }
        : t
    )
    this.setData({ trashCans })
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
