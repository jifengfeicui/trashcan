const { getAdminTrashCans, updateAdminTrashCan, deleteAdminTrashCan, getAdminTags, deleteAdminTag, createAdminTag } = require("../../utils/api/trashcan")
const { getCurrentUser } = require("../../utils/api/user")
const { isAuthenticated } = require("../../utils/auth")

Page({
  data: {
    isAdmin: false,
    currentTab: 'place',
    loading: false,
    list: [],
    page: 1,
    pageSize: 20,
    total: 0,
    totalPages: 1,
    keyword: "",
    showEditModal: false,
    editItem: null,
    editTags: "",
    editAddress: "",
    editDescription: "",
    tagList: [],
    tagLoading: false
  },

  onLoad() {
    this.checkAdminAndLoad()
  },

  onShow() {
    this.checkAdminAndLoad()
  },

  switchTab(e) {
    const tab = e.currentTarget.dataset.tab
    this.setData({ currentTab: tab })
    if (tab === 'place') {
      this.loadData()
    } else {
      this.loadTags()
    }
  },

  loadTags() {
    this.setData({ tagLoading: true })
    getAdminTags()
      .then((res) => {
        this.setData({
          tagList: res.data || []
        })
      })
      .catch((err) => {
        wx.showToast({
          title: err.message || "加载失败",
          icon: "none"
        })
      })
      .finally(() => {
        this.setData({ tagLoading: false })
      })
  },

  deleteTag(e) {
    const id = Number(e.currentTarget.dataset.id)
    const name = e.currentTarget.dataset.name
    wx.showModal({
      title: "确认删除",
      content: `确定要删除标签"${name}"吗？`,
      success: (res) => {
        if (res.confirm) {
          wx.showLoading({ title: "删除中..." })
          deleteAdminTag(id)
            .then(() => {
              wx.showToast({ title: "删除成功", icon: "success" })
              this.loadTags()
            })
            .catch((err) => {
              wx.showToast({
                title: err.message || "删除失败",
                icon: "none"
              })
            })
            .finally(() => {
              wx.hideLoading()
            })
        }
      }
    })
  },

  createTag() {
    wx.showModal({
      title: '新增标签',
      placeholderText: '请输入标签名',
      editable: true,
      success: (res) => {
        if (res.confirm && res.content && res.content.trim()) {
          const name = res.content.trim()
          wx.showLoading({ title: "创建中..." })
          createAdminTag(name)
            .then(() => {
              wx.showToast({ title: "创建成功", icon: "success" })
              this.loadTags()
            })
            .catch((err) => {
              wx.showToast({
                title: err.message || "创建失败",
                icon: "none"
              })
            })
            .finally(() => {
              wx.hideLoading()
            })
        }
      }
    })
  },

  checkAdminAndLoad() {
    if (!isAuthenticated()) {
      wx.showModal({
        title: "需要登录",
        content: "请先登录",
        success: (res) => {
          if (res.confirm) {
            wx.navigateTo({
              url: "/pages/login/index?redirect=%2Fpages%2Fadmin%2Findex"
            })
          }
        }
      })
      return
    }

    getCurrentUser()
      .then((res) => {
        const user = res.data || {}
        if (!user.is_admin) {
          wx.showModal({
            title: "权限不足",
            content: "只有管理员才能访问此页面",
            showCancel: false,
            success: () => {
              wx.navigateBack()
            }
          })
          return
        }
        this.setData({ isAdmin: true })
        if (this.data.currentTab === 'place') {
          this.loadData()
        } else {
          this.loadTags()
        }
      })
      .catch(() => {
        wx.showModal({
          title: "权限不足",
          content: "只有管理员才能访问此页面",
          showCancel: false,
          success: () => {
            wx.navigateBack()
          }
        })
      })
  },

  loadData() {
    this.setData({ loading: true })
    const { page, pageSize, keyword } = this.data

    getAdminTrashCans(page, pageSize, keyword)
      .then((res) => {
        const data = res.data || {}
        this.setData({
          list: data.list || [],
          total: data.total || 0,
          totalPages: data.total_pages || 1
        })
      })
      .catch((err) => {
        wx.showToast({
          title: err.message || "加载失败",
          icon: "none"
        })
      })
      .finally(() => {
        this.setData({ loading: false })
      })
  },

  onPullDownRefresh() {
    this.setData({ page: 1 })
    this.loadData().finally(() => wx.stopPullDownRefresh())
  },

  onReachBottom() {
    const { page, totalPages } = this.data
    if (page < totalPages) {
      this.setData({ page: page + 1 })
      this.loadData()
    }
  },

  handleKeywordInput(e) {
    this.setData({ keyword: e.detail.value })
  },

  search() {
    this.setData({ page: 1 })
    this.loadData()
  },

  openEditModal(e) {
    const item = e.currentTarget.dataset.item
    let tagsArray = item.tags_array || []
    if (typeof tagsArray === 'string') {
      try {
        tagsArray = JSON.parse(tagsArray)
      } catch (e) {
        tagsArray = []
      }
    }
    this.setData({
      showEditModal: true,
      editItem: item,
      editTags: tagsArray.join(", "),
      editAddress: item.address || "",
      editDescription: item.description || ""
    })
  },

  closeEditModal() {
    this.setData({
      showEditModal: false,
      editItem: null
    })
  },

  handleInput(e) {
    const field = e.currentTarget.dataset.field
    this.setData({
      [field]: e.detail.value
    })
  },

  saveEdit() {
    const { editItem, editTags, editAddress, editDescription } = this.data
    if (!editItem) return

    wx.showLoading({ title: "保存中..." })

    updateAdminTrashCan(editItem.id, {
      tags: editTags,
      address: editAddress,
      description: editDescription
    })
      .then(() => {
        wx.showToast({ title: "保存成功", icon: "success" })
        this.closeEditModal()
        this.loadData()
      })
      .catch((err) => {
        wx.showToast({
          title: err.message || "保存失败",
          icon: "none"
        })
      })
      .finally(() => {
        wx.hideLoading()
      })
  },

  deleteItem(e) {
    const item = e.currentTarget.dataset.item
    wx.showModal({
      title: "确认删除",
      content: `确定要删除这条记录吗？`,
      success: (res) => {
        if (res.confirm) {
          wx.showLoading({ title: "删除中..." })
          deleteAdminTrashCan(item.id)
            .then(() => {
              wx.showToast({ title: "删除成功", icon: "success" })
              this.loadData()
            })
            .catch((err) => {
              wx.showToast({
                title: err.message || "删除失败",
                icon: "none"
              })
            })
            .finally(() => {
              wx.hideLoading()
            })
        }
      }
    })
  },

  goHome() {
    wx.switchTab({
      url: "/pages/home/index"
    })
  }
})