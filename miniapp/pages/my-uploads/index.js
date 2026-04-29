const { getUserTrashCans, deleteTrashCan, updateTrashCan, uploadImage } = require("../../utils/api/trashcan")
const { resolveImageURL } = require("../../utils/request")
const { isAuthenticated } = require("../../utils/auth")

Page({
  data: {
    loadingList: false,
    trashCans: [],
    page: 1,
    pageSize: 10,
    total: 0,
    totalPages: 1,
    editingItem: null,
    editItemForm: { address: "", description: "" },
    editItemImages: [],
    editItemHasNewImage: false
  },

  onShow() {
    if (!isAuthenticated()) {
      wx.navigateBack()
      return
    }
    this.loadMyTrashCans(1)
  },

  loadMyTrashCans(page) {
    this.setData({ loadingList: true })
    getUserTrashCans(page, this.data.pageSize)
      .then((res) => {
        const payload = res.data || {}
        const list = (payload.list || []).map((item) => ({
          ...item,
          image_urls: [
            resolveImageURL(item.image_url),
            resolveImageURL(item.image_url_2),
            resolveImageURL(item.image_url_3)
          ].filter(Boolean)
        }))
        this.setData({
          trashCans: list,
          page: Number(payload.page || page),
          total: Number(payload.total || 0),
          totalPages: Number(payload.total_pages || 1)
        })
      })
      .catch((err) => wx.showToast({ title: err.message || "加载失败", icon: "none" }))
      .finally(() => this.setData({ loadingList: false }))
  },

  loadPrev() {
    if (this.data.page > 1) this.loadMyTrashCans(this.data.page - 1)
  },

  loadNext() {
    if (this.data.page < this.data.totalPages) this.loadMyTrashCans(this.data.page + 1)
  },

  handleDelete(e) {
    const id = Number(e.currentTarget.dataset.id)
    wx.showModal({
      title: "确认删除",
      content: "删除后不可恢复，是否继续？",
      success: (res) => {
        if (!res.confirm) return
        deleteTrashCan(id)
          .then(() => {
            wx.showToast({ title: "删除成功", icon: "success" })
            this.loadMyTrashCans(this.data.page)
          })
          .catch((err) => wx.showToast({ title: err.message || "删除失败", icon: "none" }))
      }
    })
  },

  startEditItem(e) {
    const id = Number(e.currentTarget.dataset.id)
    const item = this.data.trashCans.find(it => Number(it.id) === id)
    if (!item) return
    const images = [item.image_url, item.image_url_2, item.image_url_3].filter(Boolean)
    this.setData({
      editingItem: id,
      editItemForm: { address: item.address || "", description: item.description || "" },
      editItemImages: images,
      editItemHasNewImage: false
    })
  },

  cancelEditItem() {
    this.setData({ editingItem: null, editItemForm: { address: "", description: "" }, editItemImages: [], editItemHasNewImage: false })
  },

  handleItemInput(e) {
    this.setData({ [`editItemForm.${e.currentTarget.dataset.field}`]: e.detail.value })
  },

  chooseEditImage() {
    const remaining = 3 - this.data.editItemImages.length
    if (remaining <= 0) { wx.showToast({ title: "最多3张图片", icon: "none" }); return }
    wx.chooseMedia({
      count: remaining, mediaType: ["image"], sourceType: ["album", "camera"],
      success: (res) => {
        const newPaths = res.tempFiles.map(f => f.tempFilePath).slice(0, remaining)
        this.setData({ editItemImages: [...this.data.editItemImages, ...newPaths].slice(0, 3), editItemHasNewImage: true })
      }
    })
  },

  removeEditImage(e) {
    const images = [...this.data.editItemImages]
    images.splice(Number(e.currentTarget.dataset.index), 1)
    this.setData({ editItemImages: images, editItemHasNewImage: true })
  },

  saveEditItem() {
    const { editingItem, editItemForm, editItemImages, editItemHasNewImage } = this.data
    if (!editingItem) return

    const isTempPath = (p) => p && (p.startsWith('http://tmp/') || p.startsWith('https://wxfile/'))
    const buildData = (imgs) => ({
      address: editItemForm.address, description: editItemForm.description,
      image: imgs[0] || "", image_2: imgs[1] || "", image_3: imgs[2] || ""
    })

    const doUpdate = (imgs) => {
      this.setData({ loadingList: true })
      updateTrashCan(editingItem, buildData(imgs))
        .then(() => { wx.showToast({ title: "保存成功", icon: "success" }); this.setData({ editingItem: null }); this.loadMyTrashCans(this.data.page) })
        .catch((err) => wx.showToast({ title: err.message || "保存失败", icon: "none" }))
        .finally(() => this.setData({ loadingList: false }))
    }

    const tempImages = editItemImages.filter(isTempPath)
    if (!editItemHasNewImage || tempImages.length === 0) { doUpdate(editItemImages); return }

    this.setData({ loadingList: true })
    Promise.all(tempImages.map(p => uploadImage(p)))
      .then(uploaded => {
        const final = editItemImages.map(img => isTempPath(img) ? uploaded[tempImages.indexOf(img)] || img : img)
        doUpdate(final)
      })
      .catch((err) => { wx.showToast({ title: err.message || "上传图片失败", icon: "none" }); this.setData({ loadingList: false }) })
  }
})
