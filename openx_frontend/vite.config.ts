import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from "path"
import styleImport from 'vite-plugin-style-import'
import fs from 'fs'

const options = {
  key: fs.readFileSync('cert/key.pem'),
  cert: fs.readFileSync('cert/cert.pem')
}
// https://vitejs.dev/config/
export default defineConfig({
  base: "./",
  resolve: {
    alias: {
      "@": resolve(__dirname, "./src"),
      "@p": resolve(__dirname, "./proto"),
      "@pub": resolve(__dirname, "./public")
    }
  },
  plugins: [vue(
    {
      script: {
        refSugar: true
      }
    }
  ),
  // styleImport({
  //   libs: [{
  //     libraryName: 'element-plus',
  //     esModule: true,
  //     ensureStyleFile: true,
  //     resolveStyle: (name) => {
  //       name = name.slice(3)
  //       return `element-plus/packages/theme-chalk/src/${name}.scss`;
  //     },
  //     resolveComponent: (name) => {
  //       return `element-plus/lib/${name}`;
  //     },
  //   }]
  // })
  ],
  server: {
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8080',
        changeOrigin: true,
        ws: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      }
    },
    host: '0.0.0.0',
    open: true,
    port: 3001,
    https: options
  },
  build: {
    rollupOptions: {
      external: [
        
      ],
    },
  },
  optimizeDeps: {
    include: [
      'element-plus/lib/locale/lang/zh-cn',
      'element-plus/lib/locale/lang/en'
    ]
  },
})
