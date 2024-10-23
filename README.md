# Trilium Blog

Trilium Blog 是一个基于 Trilium Notes 的轻量级博客系统。它允许您将 Trilium Notes 中的笔记轻松转换为公开的博客文章，并提供了现代化的前端界面。

## 特性

- 与 Trilium Notes 无缝集成
- 支持分页的文章列表
- 单独的文章页面
- 支持代码高亮（使用 Prism.js）
- 内置评论系统（使用 Artalk）
- 响应式设计，适配移动设备
- 自动生成 sitemap.xml 和 robots.txt
- 可配置的博客名称和域名
- 暗黑模式支持
- 文章目录导航
- 图片灯箱效果（使用 Fancybox）

## 技术栈

### 后端

- Go 1.23+
- Gin Web Framework
- Redis (用于缓存)

### 前端

- Vue 3
- Vite
- Element Plus
- Axios
- Pinia
- Vue Router
- Artalk
- Fancybox
- Prism.js

## 前置要求

- Docker
- Docker Compose
- Go 1.23+ (用于开发)

## 安装和运行

1. 克隆仓库：

   ```
   git clone https://github.com/harveyTon/trilium-blog.git
   cd trilium-blog
   ```

2. 复制并编辑配置文件：

   ```
   cp ./backend/config-example.json ./backend/config.json
   ```

   编辑 `config.json`，填入您的 Trilium API URL、Token 和其他设置。

3. 使用 Docker Compose 构建和启动项目：

   ```
   docker-compose up -d --build
   ```

4. 访问 `http://localhost:8080`（或您配置的其他端口）来查看博客。

## 配置

1. 在 `config.json` 中配置您的博客信息、Trilium API 设置等。
2. 在 `frontend/src/views/Article.vue` 中配置 Artalk 评论系统。
3. 根据需要修改 `docker-compose.yml` 文件以适应您的环境。

## 使用

1. 在 Trilium Notes 中，为要发布为博客文章的笔记添加 `#blog=true` 属性。
2. 访问您配置的博客地址以查看发布的文章。

## 开发

如果您想在本地进行开发：

### 后端

1. 安装 Go 1.23+
2. 运行 `go mod tidy` 安装依赖
3. 运行 `go run main.go` 启动后端服务器

### 前端

1. 进入 `frontend` 目录
2. 运行 `npm install` 安装依赖
3. 运行 `npm run dev` 启动开发服务器

## 自定义

- 修改 `frontend/src/components/` 目录中的 Vue 组件以自定义页面布局和样式。
- 调整 `config.json` 中的设置以更改博客名称、域名等。

## 贡献

欢迎提交 Pull Requests 来改进这个项目。对于重大更改，请先开一个 issue 讨论您想要改变的内容。

## 许可证

本项目采用 MIT 许可证。详情请见 [LICENSE](LICENSE) 文件。

## 致谢

- [TriliumNext Notes](https://github.com/TriliumNext/Notes)
- [Trilium Notes](https://github.com/zadam/trilium)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Vue.js](https://vuejs.org/)
- [Element Plus](https://element-plus.org/)
- [Artalk](https://github.com/ArtalkJS/Artalk)
- [Fancybox](https://fancyapps.com/fancybox/)
- [Prism.js](https://prismjs.com/)
