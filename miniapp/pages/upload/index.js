const { createTrashCan, createTrashCanWithImages, getNearbyTrashCans, getAllTags, createAdminTag } = require("../../utils/api/trashcan")
const { getAddressByLocation } = require("../../utils/geocoder")
const { isAuthenticated } = require("../../utils/auth")

Page({
  data: {
    locationMode: "current",
    locating: false,
    loadingAddress: false,
    submitting: false,
    imagePaths: [],
    maxImages: 3,
    form: {
      latitude: "",
      longitude: "",
      address: "",
      description: ""
    },
    allTags: [],
    tagOptions: ['无'],
    selectedTagIndex: 0
  },

  onShow() {
    if (!isAuthenticated()) {
      this.goLogin()
    }
    this.loadTags()
  },

  loadTags() {
    getAllTags()
      .then((res) => {
        const tags = res.data || []
        const tagOptions = ['无', ...tags.map(t => t.name)]
        this.setData({
          allTags: tags,
          tagOptions: tagOptions
        })
      })
      .catch(() => {})
  },

  onTagChange(e) {
    const index = Number(e.detail.value)
    this.setData({ selectedTagIndex: index })
  },

  addNewTag() {
    wx.showModal({
      title: '新增标签',
      placeholderText: '请输入标签名',
      editable: true,
      success: (res) => {
        if (res.confirm && res.content && res.content.trim()) {
          const newTagName = res.content.trim()
          wx.showLoading({ title: '创建中...' })
          createAdminTag(newTagName)
            .then(() => {
              wx.showToast({ title: '创建成功', icon: 'success' })
              return this.loadTags()
            })
            .then(() => {
              const { allTags } = this.data
              const idx = allTags.findIndex(t => t.name === newTagName)
              if (idx >= 0) {
                this.setData({ selectedTagIndex: idx + 1 })
              }
            })
            .catch((err) => {
              wx.showToast({ title: err.message || '创建失败', icon: 'none' })
            })
            .finally(() => {
              wx.hideLoading()
            })
        }
      }
    })
  },

  onModeChange(e) {
    this.setData({
      locationMode: e.detail.value
    })
  },

  handleInput(e) {
    const field = e.currentTarget.dataset.field
    if (!field) {
      return
    }

    this.setData({
      [`form.${field}`]: e.detail.value
    })
  },

  getCurrentLocation() {
    this.setData({ locating: true })

    wx.getLocation({
      type: "gcj02",
      success: (res) => {
        const latitude = String(Number(res.latitude).toFixed(6))
        const longitude = String(Number(res.longitude).toFixed(6))

        this.setData({
          "form.latitude": latitude,
          "form.longitude": longitude
        })

        this.fillAddressByLocation(longitude, latitude)
      },
      fail: () => {
        wx.showToast({
          title: "定位失败",
          icon: "none"
        })
      },
      complete: () => {
        this.setData({ locating: false })
      }
    })
  },

  chooseLocationByMap() {
    wx.chooseLocation({
      success: (res) => {
        const latitude = String(Number(res.latitude).toFixed(6))
        const longitude = String(Number(res.longitude).toFixed(6))

        this.setData({
          "form.latitude": latitude,
          "form.longitude": longitude,
          "form.address": res.address || res.name || ""
        })
      },
      fail: () => {
        wx.showToast({
          title: "选点已取消",
          icon: "none"
        })
      }
    })
  },

  fillAddressByLocation(lng, lat) {
    if (!lng || !lat) {
      return
    }

    this.setData({ loadingAddress: true })

    getAddressByLocation(lng, lat)
      .then((address) => {
        if (address) {
          this.setData({
            "form.address": address
          })
        }
      })
      .catch(() => {})
      .finally(() => {
        this.setData({ loadingAddress: false })
      })
  },

  handleCoordBlur() {
    const { latitude, longitude } = this.data.form
    if (latitude && longitude) {
      this.fillAddressByLocation(longitude, latitude)
    }
  },

  chooseImage() {
    const { imagePaths, maxImages } = this.data
    const remaining = maxImages - imagePaths.length
    
    if (remaining <= 0) {
      wx.showToast({
        title: `最多上传${maxImages}张图片`,
        icon: "none"
      })
      return
    }

    wx.chooseMedia({
      count: remaining,
      mediaType: ["image"],
      sourceType: ["album", "camera"],
      success: (res) => {
        const files = res.tempFiles || []
        if (files.length === 0) {
          return
        }

        let validFiles = []
        for (const file of files) {
          if (file.size > 5 * 1024 * 1024) {
            wx.showToast({
              title: "图片不能超过5MB",
              icon: "none"
            })
            continue
          }
          validFiles.push(file.tempFilePath)
        }

        const newPaths = [...imagePaths, ...validFiles].slice(0, maxImages)
        this.setData({ imagePaths: newPaths })
      }
    })
  },

  removeImage(e) {
    const index = Number(e.currentTarget.dataset.index)
    const { imagePaths } = this.data
    imagePaths.splice(index, 1)
    this.setData({ imagePaths })
  },

  async submit() {
    if (this.data.submitting) {
      return
    }

    const that = this
    const token = wx.getStorageSync('token')
    console.log('token:', token)
    if (!token) {
      wx.showModal({
        title: "需要登录",
        content: "上传功能需要登录，是否前往登录页？",
        success: function(res) {
          if (res.confirm) {
            wx.navigateTo({
              url: "/pages/login/index?redirect=%2Fpages%2Fupload%2Findex"
            })
          }
        }
      })
      return
    }

    const { form, imagePaths, selectedTagIndex, allTags } = this.data
    const latitude = Number(form.latitude)
    const longitude = Number(form.longitude)

    if (!latitude || !longitude) {
      wx.showToast({
        title: "请先选择坐标",
        icon: "none"
      })
      return
    }

    let tagId = ''
    if (selectedTagIndex > 0) {
      tagId = allTags[selectedTagIndex - 1]?.id || ''
    }

    const payload = {
      latitude: String(latitude),
      longitude: String(longitude),
      address: form.address || "",
      description: form.description || "",
      tag_id: tagId
    }

    this.setData({ submitting: true })

    try {
      let nearbyResponse
      try {
        nearbyResponse = await getNearbyTrashCans(latitude, longitude, 0.01, 1)
      } catch (err) {
        wx.showToast({
          title: "重复检测失败，请重试",
          icon: "none"
        })
        return
      }

      const nearbyTrashCan = (nearbyResponse.data || [])[0]
      if (nearbyTrashCan) {
        const distance = Math.round(Number(nearbyTrashCan.distance || 0) * 1000)
        const address = nearbyTrashCan.address ? `\n位置：${nearbyTrashCan.address}` : ""
        const shouldContinue = await new Promise((resolve) => {
          wx.showModal({
            title: "发现附近垃圾桶",
            content: `10米内已有垃圾桶，距离约${distance}米。可能是同一个垃圾桶。${address}`,
            cancelText: "取消上传",
            confirmText: "继续上传",
            success: (res) => resolve(res.confirm),
            fail: () => resolve(false)
          })
        })

        if (!shouldContinue) {
          return
        }
      }

      console.log('going to upload, imagePaths.length:', imagePaths.length)

      if (imagePaths.length > 0) {
        await createTrashCanWithImages(payload, imagePaths)
      } else {
        await createTrashCan(payload)
      }

      wx.showToast({
        title: "上传成功",
        icon: "success"
      })
      this.resetForm()
      setTimeout(() => {
        wx.switchTab({ url: "/pages/home/index" })
      }, 1000)
    } catch (err) {
      wx.showToast({
        title: err.message || "上传失败",
        icon: "none"
      })
    } finally {
      this.setData({ submitting: false })
    }
  },

  resetForm() {
    this.setData({
      imagePaths: [],
      form: {
        latitude: "",
        longitude: "",
        address: "",
        description: "",
        tags: ""
      }
    })
  },

  goLogin() {
    wx.showModal({
      title: "需要登录",
      content: "上传功能需要登录，是否前往登录页？",
      success: (res) => {
        if (res.confirm) {
          wx.navigateTo({
            url: "/pages/login/index?redirect=%2Fpages%2Fupload%2Findex"
          })
        }
      }
    })
  },

  goHome() {
    wx.reLaunch({
      url: "/pages/home/index"
    })
  }
})
