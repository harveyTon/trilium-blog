# Trilium Blog

Trilium Blog 是一个基于 Trilium Notes 的轻量级博客系统。它允许您将 Trilium Notes 中的笔记轻松转换为公开的博客文章。

## 特性

- 与 Trilium Notes 无缝集成
- 支持分页的文章列表
- 单独的文章页面
- 支持代码高亮
- 内置评论系统（使用 Artalk）
- 响应式设计，适配移动设备
- 自动生成 sitemap.xml 和 robots.txt
- 可配置的博客名称和域名

## 安装

1. 确保您已安装 Go 1.23 或更高版本。
2. 克隆此仓库：
   ```
   git clone https://github.com/yourusername/trilium-blog.git
   cd trilium-blog
   ```
3. 安装依赖：
   ```
   go mod tidy
   ```

## 配置

1. 复制 `config.json.example` 到 `config.json`：
   ```
   cp config.json.example config.json
   ```
2. 编辑 `config.json`，填入您的 Trilium API URL、Token 和其他设置：
   ```json
   {
     "triliumApiUrl": "your_trilium_api_url",
     "triliumToken": "your_trilium_token",
     "articlesPerPage": 10,
     "blogName": "Your Blog Name",
     "domain": "https://your-domain.com"
   }
   ```

## 运行

1. 构建项目：
   ```
   go build
   ```
2. 运行服务器：
   ```
   ./trilium-blog
   ```
   默认情况下，服务器将在 `http://localhost:8080` 上运行。

## 使用

1. 在 Trilium Notes 中，为要发布为博客文章的笔记添加 `#blog=true` 属性。
2. 访问 `http://localhost:8080`（或您配置的域名）以查看您的博客。

## 自定义

- 修改 `templates/` 目录中的模板文件以自定义页面布局和样式。
- 调整 `config.json` 中的设置以更改博客名称、域名等。

## 贡献

欢迎提交 Pull Requests 来改进这个项目。对于重大更改，请先开一个 issue 讨论您想要改变的内容。

## 许可证

本项目采用 MIT 许可证。详情请见 [LICENSE](LICENSE) 文件。

## 致谢
- [TriliumNext Notes](https://github.com/TriliumNext/Notes)
- [Trilium Notes](https://github.com/zadam/trilium)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Artalk](https://github.com/ArtalkJS/Artalk)
- [highlight.js](https://highlightjs.org/)
