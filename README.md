# Trilium Blog

中文 | **[English](README.en.md)**

**演示：[虎笺.com](https://虎笺.com/)**

基于 [TriliumNext Notes](https://github.com/TriliumNext/Trilium) 的轻量级博客系统。将 Trilium 中标记为 `#blog=true` 的笔记发布为博客文章。

## 特性

- 通过 ETAPI 与 Trilium Notes 集成，自动读取 `#blog=true` 笔记
- 首页包含精选文章、最新文章分页列表与居中的全局搜索框
- 独立搜索页与搜索预览面板
- 精选文章单卡片全宽轮播，支持左右箭头、圆点与触摸滑动切换
- AI summary / code summary 双摘要体系
- 文章页优先加载正文，AI summary 异步生成与轮询回填，不阻塞阅读
- 阅读模式：沉浸式排版、阅读进度、目录抽屉、宽度/密度/字体/主题设置
- 文章目录（TOC）、阅读进度条、代码复制按钮、代码语言标签
- 后端使用 Chroma 识别代码语言，前端使用 Shiki 渲染高亮
- 图片灯箱（Fancybox）
- 外部图片代理与内置 `/api/imageproxy`
- Redis 缓存、启动时异步预加载文章内容、SQLite 摘要存储
- 暗黑模式与移动端适配
- `sitemap.xml` / `robots.txt`

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

默认部署方式就是 Docker Compose。当前仓库的 `docker-compose.yml` 会同时启动：

- `trilium-blog`
- `redis`

其中：

- Redis 地址固定为容器内 `redis:6379`，不需要额外配置。
- AI / code summary 的 SQLite 数据库位于容器内 `/app/data/summaries.db`。
- `trilium-blog-data` volume 会持久化摘要数据库，避免 `docker compose up -d --build` 后重复全量生成 AI summary。

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
| `DOMAIN` | 否 | — | 博客域名，用于 sitemap、页面链接与内部资源判断 |
| `ARTICLES_PER_PAGE` | 否 | `9` | 每页文章数 |
| `PORT` | 否 | `8080` | 服务端口 |
| `LOCALE` | 否 | `zh-CN` | 博客语言，支持 `zh-CN`（中文）和 `en`（英文） |
| `ADMIN_TOKEN` | 否 | — | 管理页面令牌，设置后启用 `/admin` 缓存管理页面 |
| `LOG_LEVEL` | 否 | `info` | 日志级别：`debug`、`info`、`warn`、`error`、`fatal` |
| `IMAGE_PROXY_ENABLED` | 否 | `false` | 启用外部图片代理 |
| `IMAGE_PROXY_BASE_URL` | 否 | — | 外部图片代理 URL（留空则使用内置 `/api/imageproxy`） |
| `AI_SUMMARY_ENABLED` | 否 | `false` | 开启摘要子系统 |
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
- 一旦 AI 摘要就绪，文章页会自动更新显示；列表页和精选文章若使用的是 AI 摘要，会带有 `AI` 标识。
- AI 生成上下文会同时包含文章标题和正文，而不是只使用正文。
- `posts` 接口已直接返回 `summaries`，前端会优先复用，不会在已有结果时重复请求摘要接口。

如果只想保留本地摘要、不发起 AI 请求，可以设置：

```env
AI_SUMMARY_ENABLED=false
```

或

```env
AI_SUMMARY_MODE=code
```

## 缓存与预加载

- 所有文章列表、内容、附件均通过 Redis 缓存（策略驱动的 TTL 管理）。
- 服务启动后会在后台异步预加载全部 `#blog=true` 文章的原始内容到 Redis，首次访问时直接命中缓存，无需等待 Trilium ETAPI 响应。
- 预加载仅缓存原始内容，不触发 code summary 或 AI summary 生成。
- 如果 Redis 不可用，自动降级为无缓存模式，所有请求直接转发至 Trilium。

### 缓存管理

设置 `ADMIN_TOKEN` 后，访问 `/admin` 进入缓存管理页面（支持中英双语）：

- 查看 Redis 连接状态与各缓存类型的 key 数量、TTL 信息
- 按类型或全局清除缓存
- 按 note ID 或 attachment ID 精确失效
- 手动触发预加载（不会在清除缓存后自动触发）

## 首页与文章页行为

### 首页

- 首页由“精选文章”与“最新文章”两部分组成。
- 精选文章来自 `#blogtop=true` 标签。
- 最新文章支持分页，分页会同步更新浏览器 URL。
- 第 1 页页面标题为 `BLOG_TITLE | BLOG_SUBTITLE`。
- 第 2 页及之后页面标题会追加页码，例如：`BLOG_TITLE | BLOG_SUBTITLE - 第 2 页`。

### 文章页

- 文章内容、TOC、摘要、代码块信息均由后端返回。
- 文章页标题格式为：`文章标题 - BLOG_TITLE | BLOG_SUBTITLE`。
- 正常模式下展示摘要、正文与源码链接。
- 阅读模式与普通模式共用同一路由，仅通过页面级 class 切换。
- 阅读模式偏好会持久化保存，包括：
  - TOC 收起状态
  - 宽度
  - 密度
  - 字体大小
  - 阅读主题

## 代码块与高亮

- 后端会在处理文章 HTML 时抽取代码块元数据。
- 代码语言识别按以下优先级进行：
  1. Trilium 返回的 MIME class 名（内置 120+ MIME type 映射表）
  2. [Chroma](https://github.com/alecthomas/chroma) 词法分析
  3. [enry](https://github.com/go-enry/go-enry) 统计分类器
  4. 兜底为 `plaintext`
- 前端使用 [Shiki](https://github.com/shikijs/shiki) 进行代码高亮渲染。
- 支持代码语言标签、复制按钮、行号与亮暗主题切换。

## 使用

在 Trilium Notes 中：

- 为要发布的笔记添加 `#blog=true`
- 为要加入精选文章的笔记添加 `#blogtop=true`

博客会自动读取并展示这些内容。

## 技术栈

**后端：** Go 1.25 / Gin / goquery / bluemonday / Redis / SQLite / Chroma / enry

**前端：** Vue 3 / Vite / Element Plus / Pinia / Vue Router / Fancybox / Shiki

## 致谢

- [Trilium](https://github.com/TriliumNext/Trilium)
- [Trilium Notes](https://github.com/zadam/trilium)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Vue.js](https://vuejs.org/)
- [Element Plus](https://element-plus.org/)
- [Fancybox](https://fancyapps.com/fancybox/)
- [Shiki](https://github.com/shikijs/shiki)
- [Chroma](https://github.com/alecthomas/chroma)
- [enry](https://github.com/go-enry/go-enry)

## 许可证

[MIT](LICENSE)
