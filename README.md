# Trilium Blog

基于 [TriliumNext Notes](https://github.com/TriliumNext/Trilium) 的轻量级博客系统。将 Trilium 中标记为 `#blog=true` 的笔记发布为博客文章。

## 特性

- 与 Trilium Notes 无缝集成（通过 ETAPI）
- 首页精选文章、最新文章分页与顶部搜索框
- 独立搜索页与搜索预览
- 代码高亮（highlight.js，按需引入）
- 文章目录导航（TOC）与阅读进度
- 图片灯箱（Fancybox）
- AI / code summary 双层摘要能力（可选，异步回填）
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

如果启用了 AI summary，Docker 配置会自动持久化摘要数据库，避免 `docker compose up -d --build` 后重复全量生成。

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
| `BLOG_TITLE` | 否 | — | 博客主标题 |
| `BLOG_SUBTITLE` | 否 | — | 博客副标题 |
| `DOMAIN` | 否 | — | 博客域名（用于内部资产判断） |
| `ARTICLES_PER_PAGE` | 否 | `9` | 每页文章数 |
| `PORT` | 否 | `8080` | 服务端口 |
| `IMAGE_PROXY_ENABLED` | 否 | `false` | 启用外部图片代理 |
| `IMAGE_PROXY_BASE_URL` | 否 | — | 外部图片代理 URL（留空则使用内置 `/api/imageproxy`） |
| `AI_SUMMARY_ENABLED` | 否 | `false` | 启用 AI 摘要子系统 |
| `AI_SUMMARY_PROVIDER` | 否 | `openai-compatible` | AI provider 类型 |
| `AI_SUMMARY_BASE_URL` | 否 | — | AI 接口基础地址 |
| `AI_SUMMARY_API_KEY` | 否 | — | AI 接口密钥 |
| `AI_SUMMARY_MODEL` | 否 | — | AI 模型名 |
| `AI_SUMMARY_PROMPT` | 否 | 内置默认值 | AI 摘要系统提示词 |
| `AI_SUMMARY_MODE` | 否 | `code` | `code` 仅生成本地摘要，`ai` 在保留 code summary 的同时异步生成 AI summary |
| `AI_SUMMARY_CONCURRENCY` | 否 | `2` | AI 摘要并发 worker 数 |
| `AI_SUMMARY_RATE_LIMIT_MS` | 否 | `1200` | AI 摘要请求间隔（毫秒） |
| `AI_SUMMARY_TIMEOUT_MS` | 否 | `60000` | 单次 AI 请求超时（毫秒） |
| `AI_SUMMARY_MAX_INPUT_CHARS` | 否 | `12000` | 发送给 AI 的正文最大字符数 |

AI summary 的 SQLite 文件由后端内部管理；在 Docker 环境下默认位于容器内 `/app/data/summaries.db`，并通过 `docker-compose.yml` 中的 volume 持久化，不需要用户配置路径。

## AI 摘要

项目内置两层摘要能力：

- `code summary`：同步生成、实时落库，用于保证列表页和文章页始终有可展示的摘要。
- `ai summary`：异步生成并回填数据库，不会阻塞文章正文加载。

当前 AI 行为如下：

- 文章详情页会优先正常加载正文。
- 如果启用了 AI summary 且当前文章还没有 AI 结果，后端会立即排队生成任务。
- 前端会显示 AI 摘要卡片的“生成中”状态，并轮询等待结果回填。
- 一旦 AI 摘要就绪，文章页会自动更新显示；列表页若使用的是 AI 摘要，会带有明显的 `AI` 标识。
- AI 生成上下文会同时包含文章标题和正文，而不是只使用正文。

如果只想保留本地摘要、不发起 AI 请求，可以设置：

```env
AI_SUMMARY_ENABLED=false
```

或

```env
AI_SUMMARY_MODE=code
```

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
