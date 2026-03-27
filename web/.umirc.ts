import { defineConfig } from "umi";

export default defineConfig({
  routes: [
    { path: "/", component: "index" },
    { path: "/docs", component: "docs" },
  ],
  npmClient: 'pnpm',
  proxy: {
    '/proxy/': {
      target: 'http://127.0.0.1:8182',
      changeOrigin: true,
      pathRewrite: { '^/proxy/': '/' },
    },
  },
});
