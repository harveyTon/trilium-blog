# Trilium Blog

基于 [TriliumNext Notes](https://github.com/TriliumNext/Trilium) 的轻量级博客系统。将 Trilium 中标记为 `#blog=true` 的笔记发布为博客文章。

## 特性

- 与 Trilium Notes 无缝集成（通过 ETAPI）
- 分页文章列表 + 文章详情页
- 代码高亮（highlight.js，按需引入）
- 文章目录导航（TOC）
- 图片灯箱（Fancybox）
- 可选图片代理（外部代理或内置 fallback）
- 内置评论系统（Artalk，可配置）
- 暗黑模式（Dark / Light 同等支持）
- 响应式设计
- sitemap.xml / robots.txt

## 快速开始

### Docker Compose（推荐）

```bash
git clone https://github.com/harveyTon/trilium-blog.git
cd trilium-blog

# 创建配置文件
cp .env.example .env
# 编辑 .env，填入 TRILIUM_API_URL、TRILIUM_TOKEN 等
vim .env

docker compose up -d --build
```

访问 `http://localhost:8080`（端口可通过 `.env` 中的 `PORT` 修改）。

### 本地开发

**后端**（需要 Go 1.25+）：

```bash
cd backend
go mod download
# 设置环境变量或 export TRILIUM_API_URL=... TRILIUM_TOKEN=... 等
go run main.go
```

**前端**（需要 Node 24+）：

```bash
cd frontend
npm install
npm run dev
```

## 配置

所有配置通过环境变量（`.env` 文件）管理：

| 变量 | 必填 | 默认值 | 说明 |
|------|------|--------|------|
| `TRILIUM_API_URL` | 是 | — | Trilium ETAPI 地址 |
| `TRILIUM_TOKEN` | 是 | — | ETAPI Token |
| `BLOG_NAME` | 否 | — | 博客名称 |
| `BLOG_TITLE` | 否 | — | 页面标题 |
| `DOMAIN` | 否 | — | 博客域名（用于内部资产判断） |
| `ARTICLES_PER_PAGE` | 否 | `9` | 每页文章数 |
| `PORT` | 否 | `8080` | 服务端口 |
| `IMAGE_PROXY_ENABLED` | 否 | `false` | 启用外部图片代理 |
| `IMAGE_PROXY_BASE_URL` | 否 | — | 外部图片代理 URL（留空则使用内置 `/api/imageproxy`） |

## 使用

在 Trilium Notes 中，为要发布的笔记添加 `#blog=true` 属性，博客会自动展示。

## 技术栈

**后端：** Go 1.25 / Gin / goquery / bluemonday

**前端：** Vue 3 / Vite / Element Plus / Pinia / Vue Router / Artalk / Fancybox / highlight.js

## 致谢

- [Trilium](https://github.com/TriliumNext/Trilium)
- [Trilium Notes](https://github.com/zadam/trilium)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Vue.js](https://vuejs.org/)
- [Element Plus](https://element-plus.org/)
- [Artalk](https://github.com/ArtalkJS/Artalk)
- [Fancybox](https://fancyapps.com/fancybox/)
- [highlight.js](https://github.com/highlightjs/highlight.js)

## 许可证

[MIT](LICENSE)
