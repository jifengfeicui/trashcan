import { fileURLToPath, URL } from 'node:url'
import { networkInterfaces } from 'os'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import basicSsl from '@vitejs/plugin-basic-ssl'

// 获取本机局域网IP地址
function getLocalIP() {
  const interfaces = networkInterfaces()
  for (const name of Object.keys(interfaces)) {
    for (const iface of interfaces[name]) {
      // 跳过内部（即127.0.0.1）和非IPv4地址
      if (iface.family === 'IPv4' && !iface.internal) {
        return iface.address
      }
    }
  }
  return 'localhost'
}

// 获取后端API地址
function getBackendURL() {
  // 优先使用环境变量
  if (process.env.VITE_API_URL) {
    return process.env.VITE_API_URL
  }
  
  // 如果设置了VITE_USE_LOCAL_IP，使用局域网IP
  if (process.env.VITE_USE_LOCAL_IP === 'true') {
    const localIP = getLocalIP()
    return `http://${localIP}:38080`
  }
  
  // 默认使用localhost
  return 'http://localhost:38080'
}

const backendURL = getBackendURL()

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
    basicSsl(), // 自动生成自签名证书以支持HTTPS
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  server: {
    host: '0.0.0.0', // 允许局域网访问
    port: 5173,
    https: true, // 启用 HTTPS（使用basicSsl插件自动生成证书）
    proxy: {
      '/api': {
        target: backendURL,
        changeOrigin: true,
        secure: false, // 如果使用自签名证书，设置为 false
        rewrite: (path) => path.replace(/^\/api/, '/api')
      },
      '/uploads': {
        target: backendURL,
        changeOrigin: true,
        secure: false, // 如果使用自签名证书，设置为 false
      }
    }
  }
})
