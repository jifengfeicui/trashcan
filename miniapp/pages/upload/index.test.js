const assert = require("node:assert/strict")

const api = require("../../utils/api/trashcan")

let nearbyResult
let nearbyError
let modalConfirm
let nearbyArgs
let createCalls
let imageCreateCalls
let modalOptions
let toasts

api.getNearbyTrashCans = async (...args) => {
  nearbyArgs = args
  if (nearbyError) throw nearbyError
  return { data: nearbyResult }
}
api.createTrashCan = async () => { createCalls++ }
api.createTrashCanWithImages = async () => { imageCreateCalls++ }

global.wx = {
  getStorageSync: () => "token",
  showModal: (options) => {
    modalOptions = options
    options.success({ confirm: modalConfirm })
  },
  showToast: (options) => toasts.push(options),
  switchTab: () => {}
}

let page
global.Page = (options) => { page = options }
global.setTimeout = (fn) => fn()
require("./index")

page.setData = function (updates) {
  Object.assign(this.data, updates)
}

function reset(images = []) {
  nearbyResult = []
  nearbyError = null
  modalConfirm = false
  nearbyArgs = null
  createCalls = 0
  imageCreateCalls = 0
  modalOptions = null
  toasts = []
  page.data.submitting = false
  page.data.imagePaths = images
  page.data.selectedTagIndex = 0
  page.data.allTags = []
  page.data.form = {
    latitude: "30",
    longitude: "120",
    address: "测试地址",
    description: ""
  }
}

async function run() {
  reset()
  await page.submit()
  assert.deepEqual(nearbyArgs, [30, 120, 0.01, 1])
  assert.equal(createCalls, 1)
  assert.equal(modalOptions, null)

  reset()
  nearbyResult = [{ distance: 0.006, address: "附近地址" }]
  await page.submit()
  assert.equal(createCalls, 0)
  assert.match(modalOptions.content, /距离约6米/)
  assert.equal(page.data.form.latitude, "30")
  assert.equal(page.data.submitting, false)

  reset(["photo.jpg"])
  nearbyResult = [{ distance: 0.01 }]
  modalConfirm = true
  await page.submit()
  assert.equal(imageCreateCalls, 1)
  assert.equal(modalOptions.cancelText, "取消上传")
  assert.equal(modalOptions.confirmText, "继续上传")

  reset()
  nearbyError = new Error("network error")
  await page.submit()
  assert.equal(createCalls, 0)
  assert.equal(toasts.at(-1).title, "重复检测失败，请重试")
  assert.equal(page.data.submitting, false)
}

run().catch((err) => {
  console.error(err)
  process.exitCode = 1
})
