const { login, register } = require("../../utils/api/user")
const { setToken, setUserInfo, isAuthenticated } = require("../../utils/auth")

Page({
  data: {
    mode: "login",
    loading: false,
    error: "",
    redirect: "/pages/home/index",
    loginForm: {
      username: "",
      password: ""
    },
    registerForm: {
      username: "",
      password: "",
      confirmPassword: ""
    }
  },

  onLoad(options) {
    if (options && options.redirect) {
      this.setData({
        redirect: decodeURIComponent(options.redirect)
      })
    }
  },

  onShow() {
    if (isAuthenticated()) {
      this.redirectAfterAuth()
    }
  },

  switchMode(e) {
    const mode = e.currentTarget.dataset.mode
    if (!mode || mode === this.data.mode) {
      return
    }

    this.setData({
      mode,
      error: ""
    })
  },

  handleInput(e) {
    const field = e.currentTarget.dataset.field
    const mode = e.currentTarget.dataset.mode
    if (!field || !mode) {
      return
    }

    this.setData({
      [`${mode}.${field}`]: e.detail.value
    })
  },

  submit() {
    if (this.data.loading) {
      return
    }

    if (this.data.mode === "login") {
      this.handleLogin()
      return
    }

    this.handleRegister()
  },

  handleLogin() {
    const { username, password } = this.data.loginForm
    if (!username || !password) {
      this.setData({ error: "请输入用户名和密码" })
      return
    }

    this.setData({ loading: true, error: "" })

    login(username, password)
      .then((res) => {
        const payload = res.data || {}
        setToken(payload.token || "")
        setUserInfo(payload.user || null)
        wx.showToast({ title: "登录成功", icon: "success" })
        this.redirectAfterAuth()
      })
      .catch((err) => {
        this.setData({ error: err.message || "登录失败" })
      })
      .finally(() => {
        this.setData({ loading: false })
      })
  },

  handleRegister() {
    const { username, password, confirmPassword } = this.data.registerForm
    if (!username || !password || !confirmPassword) {
      this.setData({ error: "请完整填写注册信息" })
      return
    }

    if (password !== confirmPassword) {
      this.setData({ error: "两次输入的密码不一致" })
      return
    }

    if (password.length < 6) {
      this.setData({ error: "密码至少 6 位" })
      return
    }

    this.setData({ loading: true, error: "" })

    register(username, password)
      .then((res) => {
        const payload = res.data || {}
        setToken(payload.token || "")
        setUserInfo(payload.user || null)
        wx.showToast({ title: "注册成功", icon: "success" })
        this.redirectAfterAuth()
      })
      .catch((err) => {
        this.setData({ error: err.message || "注册失败" })
      })
      .finally(() => {
        this.setData({ loading: false })
      })
  },

  redirectAfterAuth() {
    const url = this.data.redirect || "/pages/home/index"
    wx.reLaunch({ url })
  },

  goHome() {
    wx.reLaunch({
      url: "/pages/home/index"
    })
  }
})
