const { createTrashCan } = require("../../utils/api/trashcan")
const { uploadFile } = require("../../utils/request")
const { getAddressByLocation } = require("../../utils/geocoder")
const { isAuthenticated } = require("../../utils/auth")

Page({
  data: {
    locationMode: "current",
    locating: false,
    loadingAddress: false,
    submitting: false,
    imagePath: "",
    form: {
      latitude: "",
      longitude: "",
      address: "",
      description: ""
    }
  },

  onShow() {
    if (!isAuthenticated()) {
      this.goLogin()
    }
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
    wx.chooseMedia({
      count: 1,
      mediaType: ["image"],
      sourceType: ["album", "camera"],
      success: (res) => {
        const file = (res.tempFiles || [])[0]
        if (!file) {
          return
        }

        if (file.size > 5 * 1024 * 1024) {
          wx.showToast({
            title: "图片不能超过5MB",
            icon: "none"
          })
          return
        }

        this.setData({ imagePath: file.tempFilePath })
      }
    })
  },

  removeImage() {
    this.setData({ imagePath: "" })
  },

  submit() {
    if (!isAuthenticated()) {
      this.goLogin()
      return
    }

    const { form, imagePath } = this.data
    const latitude = Number(form.latitude)
    const longitude = Number(form.longitude)

    if (!latitude || !longitude) {
      wx.showToast({
        title: "请先选择坐标",
        icon: "none"
      })
      return
    }

    const payload = {
      latitude: String(latitude),
      longitude: String(longitude),
      address: form.address || "",
      description: form.description || ""
    }

    this.setData({ submitting: true })

    const submitTask = imagePath
      ? uploadFile({
          url: "/trashcans",
          filePath: imagePath,
          formData: payload
        })
      : createTrashCan(payload)

    submitTask
      .then(() => {
        wx.showToast({
          title: "上传成功",
          icon: "success"
        })
        this.resetForm()
      })
      .catch((err) => {
        wx.showToast({
          title: err.message || "上传失败",
          icon: "none"
        })
      })
      .finally(() => {
        this.setData({ submitting: false })
      })
  },

  resetForm() {
    this.setData({
      imagePath: "",
      form: {
        latitude: "",
        longitude: "",
        address: "",
        description: ""
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
