// 局域网访问模式的配置文件
// 使用方式: pnpm dev:local 或 vite --config vite.config.local.js
import { fileURLToPath, URL } from 'node:url'
import { networkInterfaces } from 'os'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// 获取本机局域网IP地址
function getLocalIP() {
  const interfaces = networkInterfaces()
  for (const name of Object.keys(interfaces)) {
    for (const iface of interfaces[name]) {
      if (iface.family === 'IPv4' && !iface.internal) {
        return iface.address
      }
    }
  }
  return 'localhost'
}

const localIP = getLocalIP()
const backendURL = `http://${localIP}:38080`

console.log(`🌐 前端访问地址: http://${localIP}:5173`)
console.log(`🔗 后端API地址: ${backendURL}`)

export default defineConfig({
  plugins: [
    vue(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  server: {
    host: '0.0.0.0', // 允许局域网访问
    port: 5173,
    proxy: {
      '/api': {
        target: backendURL,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '/api')
      },
      '/uploads': {
        target: backendURL,
        changeOrigin: true,
      }
    }
  }
})

