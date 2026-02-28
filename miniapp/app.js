App({
  globalData: {
    apiOrigin: ""
  },

  onLaunch() {
    const config = require('./config/index')
    this.globalData.apiOrigin = config.API_ORIGIN
  }
})
