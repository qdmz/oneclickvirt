import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  return {
    base: mode === 'production' ? '/' : '/',
    plugins: [vue()],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      }
    },
    server: {
      port: parseInt(env.VITE_CLI_PORT) || 8080,
      open: true,
      proxy: {
        [env.VITE_BASE_API || '/api']: {
          target: `${env.VITE_BASE_PATH}:${env.VITE_SERVER_PORT}` || 'http://0.0.0.0:8888',
          changeOrigin: true,          ws: true, // 代理 WebSocket 连接（ws:// 和 wss://）          // 可选：调试日志
          configure: (proxy, options) => {
            proxy.on('proxyReq', (proxyReq, req, res) => {
              // console.log(`代理请求: ${req.method} ${req.url} -> ${options.target}${proxyReq.path}`)
            })
          }
        }
      }
    },
    css: {
      preprocessorOptions: {
        scss: {
          additionalData: `@use '@/style/theme.scss' as *;`
        }
      }
    },
    build: {
      outDir: 'dist',
      assetsDir: 'assets',
      sourcemap: false,
      minify: 'esbuild',
      rollupOptions: {
        output: {
          chunkFileNames: 'js/[name]-[hash].js',
          entryFileNames: 'js/[name]-[hash].js',
          assetFileNames: 'assets/[name]-[hash].[ext]',
          // 重型依赖分包（P1-7）：把 vue 运行时、element-plus、echarts、
          // xterm、novnc 从默认 vendor 里拆出，首屏只解析真正用到的部分，
          // 且库升级时业务 chunk 哈希不被连带污染、利于长缓存。
          manualChunks(id) {
            if (!id.includes('node_modules')) return
            if (id.includes('@novnc')) return 'vendor-novnc'
            if (id.includes('@xterm')) return 'vendor-xterm'
            if (id.includes('echarts') || id.includes('zrender')) return 'vendor-echarts'
            if (id.includes('element-plus') || id.includes('@element-plus')) return 'vendor-element'
            if (id.includes('/vue/') || id.includes('vue-router') || id.includes('pinia') || id.includes('vue-i18n') || id.includes('@vue/')) return 'vendor-vue'
          }
        }
      },
      esbuild: {
        drop: mode === 'production' ? ['console', 'debugger'] : []
      }
    }
  }
})
