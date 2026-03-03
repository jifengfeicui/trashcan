const { getCurrentUser, updateCurrentUser } = require("../../utils/api/user")
const { getUserTrashCans, deleteTrashCan, updateTrashCan } = require("../../utils/api/trashcan")
const { resolveImageURL } = require("../../utils/request")
const { isAuthenticated, clearAuth } = require("../../utils/auth")

Page({
  data: {
    authed: false,
    loadingUser: false,
    loadingList: false,
    userInfo: null,
    editing: false,
    editForm: {
      nickname: "",
      avatar: ""
    },
    trashCans: [],
    page: 1,
    pageSize: 10,
    total: 0,
    totalPages: 1,
    editingItem: null,
    editItemForm: {
      address: "",
      description: ""
    },
    editItemImages: [],
    editItemHasNewImage: false
  },

  onShow() {
    this.bootstrap()
  },

  bootstrap() {
    const authed = isAuthenticated()
    this.setData({ authed })

    if (!authed) {
      return
    }

    this.loadCurrentUser()
    this.loadMyTrashCans(1)
  },

  loadCurrentUser() {
    this.setData({ loadingUser: true })

    getCurrentUser()
      .then((res) => {
        this.setData({
          userInfo: res.data || null
        })
      })
      .catch((err) => {
        wx.showToast({
          title: err.message || "获取用户失败",
          icon: "none"
        })
      })
      .finally(() => {
        this.setData({ loadingUser: false })
      })
  },

  loadMyTrashCans(page) {
    this.setData({ loadingList: true })

    getUserTrashCans(page, this.data.pageSize)
      .then((res) => {
        const payload = res.data || {}
        const list = (payload.list || []).map((item) => ({
          ...item,
          image_url: resolveImageURL(item.image_url)
        }))

        this.setData({
          trashCans: list,
          page: Number(payload.page || page),
          total: Number(payload.total || 0),
          totalPages: Number(payload.total_pages || 1)
        })
      })
      .catch((err) => {
        wx.showToast({
          title: err.message || "加载失败",
          icon: "none"
        })
      })
      .finally(() => {
        this.setData({ loadingList: false })
      })
  },

  loadPrev() {
    const nextPage = this.data.page - 1
    if (nextPage < 1) {
      return
    }

    this.loadMyTrashCans(nextPage)
  },

  loadNext() {
    const nextPage = this.data.page + 1
    if (nextPage > this.data.totalPages) {
      return
    }

    this.loadMyTrashCans(nextPage)
  },

  handleDelete(e) {
    const id = Number(e.currentTarget.dataset.id)
    if (!id) {
      return
    }

    wx.showModal({
      title: "确认删除",
      content: "删除后不可恢复，是否继续？",
      success: (res) => {
        if (!res.confirm) {
          return
        }

        deleteTrashCan(id)
          .then(() => {
            wx.showToast({
              title: "删除成功",
              icon: "success"
            })

            const currentPage = this.data.page
            this.loadMyTrashCans(currentPage)
          })
          .catch((err) => {
            wx.showToast({
              title: err.message || "删除失败",
              icon: "none"
            })
          })
      }
    })
  },

  logout() {
    clearAuth()
    this.setData({
      authed: false,
      userInfo: null,
      trashCans: []
    })

    wx.showToast({
      title: "已退出",
      icon: "success"
    })
  },

  goLogin() {
    wx.navigateTo({
      url: "/pages/login/index?redirect=%2Fpages%2Fprofile%2Findex"
    })
  },

  goHome() {
    wx.reLaunch({
      url: "/pages/home/index"
    })
  },

  goUpload() {
    wx.navigateTo({
      url: "/pages/upload/index"
    })
  },

  startEdit() {
    const { userInfo } = this.data
    this.setData({
      editing: true,
      editForm: {
        nickname: userInfo.nickname || "",
        avatar: userInfo.avatar || ""
      }
    })
  },

  cancelEdit() {
    this.setData({ editing: false })
  },

  handleInput(e) {
    const field = e.currentTarget.dataset.field
    const value = e.detail.value
    this.setData({
      [`editForm.${field}`]: value
    })
  },

  chooseAvatar() {
    wx.chooseMedia({
      count: 1,
      mediaType: ["image"],
      sourceType: ["album", "camera"],
      success: (res) => {
        const tempFilePath = res.tempFiles[0].tempFilePath
        this.setData({
          "editForm.avatar": tempFilePath
        })
      }
    })
  },

  saveProfile() {
    const { editForm } = this.data
    this.setData({ loadingUser: true })

    updateCurrentUser(editForm.nickname, editForm.avatar)
      .then((res) => {
        wx.showToast({
          title: "保存成功",
          icon: "success"
        })
        this.setData({
          userInfo: res.data || null,
          editing: false
        })
      })
      .catch((err) => {
        wx.showToast({
          title: err.message || "保存失败",
          icon: "none"
        })
      })
      .finally(() => {
        this.setData({ loadingUser: false })
      })
  },

  startEditItem(e) {
    const id = Number(e.currentTarget.dataset.id)
    const item = this.data.trashCans.find(it => Number(it.id) === id)
    if (!item) return

    const images = []
    if (item.image_url) images.push(item.image_url)
    if (item.image_url_2) images.push(item.image_url_2)
    if (item.image_url_3) images.push(item.image_url_3)

    this.setData({
      editingItem: id,
      editItemForm: {
        address: item.address || "",
        description: item.description || ""
      },
      editItemImages: images,
      editItemHasNewImage: false
    })
  },

  cancelEditItem() {
    this.setData({
      editingItem: null,
      editItemForm: { address: "", description: "" },
      editItemImages: [],
      editItemHasNewImage: false
    })
  },

  chooseEditImage() {
    const { editItemImages } = this.data
    const remaining = 3 - editItemImages.length

    if (remaining <= 0) {
      wx.showToast({ title: "最多3张图片", icon: "none" })
      return
    }

    wx.chooseMedia({
      count: remaining,
      mediaType: ["image"],
      sourceType: ["album", "camera"],
      success: (res) => {
        const files = res.tempFiles || []
        const newPaths = files.map(f => f.tempFilePath).slice(0, remaining)
        this.setData({
          editItemImages: [...editItemImages, ...newPaths].slice(0, 3),
          editItemHasNewImage: true
        })
      }
    })
  },

  removeEditImage(e) {
    const index = Number(e.currentTarget.dataset.index)
    const { editItemImages } = this.data
    editItemImages.splice(index, 1)
    this.setData({
      editItemImages,
      editItemHasNewImage: true
    })
  },

  handleItemInput(e) {
    const field = e.currentTarget.dataset.field
    const value = e.detail.value
    this.setData({
      [`editItemForm.${field}`]: value
    })
  },

  saveEditItem() {
    const { editingItem, editItemForm } = this.data
    if (!editingItem) return

    this.setData({ loadingList: true })

    updateTrashCan(editingItem, {
      address: editItemForm.address,
      description: editItemForm.description
    })
      .then(() => {
        wx.showToast({ title: "保存成功", icon: "success" })
        this.setData({ editingItem: null })
        this.loadMyTrashCans(this.data.page)
      })
      .catch((err) => {
        wx.showToast({ title: err.message || "保存失败", icon: "none" })
      })
      .finally(() => {
        this.setData({ loadingList: false })
      })
  }
})
