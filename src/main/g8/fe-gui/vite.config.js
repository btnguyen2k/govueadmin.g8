/** @type {import('vite').UserConfig} */
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from "path";

// https://vitejs.dev/config/
export default defineConfig(({ command, mode }) => {
  const path = require('path')
  const env = loadEnv(mode, process.cwd(), '')
  process.env = {...process.env, ...env};
  return {
    server: {
      port: 8080,
      host: "0.0.0.0",
      proxy: {
        // https://vitejs.dev/config/server-options.html
        // '/api': {
        //   target: 'http://localhost:8000/',
        //   changeOrigin: true,
        //   rewrite: (path) => path.replace(/^\/api/, '')
        // },
      },
    },
    plugins: [vue()],
    base: env.BASE_URL ?? '/app/',
    resolve: {
      extensions: ['.mjs', '.js', '.ts', '.jsx', '.tsx', '.json', '.vue', 'scss'],
      alias: [
        {
          // this is required for the SCSS modules
          find: /^~(.*)$/,
          replacement: '$1',
        },
        {
          find: '@/',
          replacement: `${path.resolve(__dirname, 'src')}/`,
        },
        {
          find: '@',
          replacement: path.resolve(__dirname, '/src'),
        },
      ],
    },
    define: {
      'process.env': process.env,
    },
  }
})
