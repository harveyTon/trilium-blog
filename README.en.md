# Trilium Blog

**[中文](README.md)** | English

**Demo: [虎笺.com](https://虎笺.com/)**

A lightweight blog system powered by [TriliumNext Notes](https://github.com/TriliumNext/Trilium). Publish notes tagged with `#blog=true` from Trilium as blog articles.

## Features

- Integrates with Trilium Notes via ETAPI, automatically reading `#blog=true` notes
- Homepage with featured posts, paginated latest posts, and a centered global search box
- Dedicated search page and search preview panel
- Featured post single-card full-width carousel with arrow, dot, and swipe navigation
- AI summary / code summary dual summary system
- Article page loads content first; AI summary is generated asynchronously and polled back, non-blocking
- Reading mode: immersive typography, reading progress, TOC drawer, width/density/font/theme settings
- Table of contents (TOC), reading progress bar, code copy button, code language labels
- Backend uses Chroma + enry for code language detection, frontend uses Shiki for syntax highlighting
- Image lightbox (Fancybox)
- External image proxy with built-in `/api/imageproxy`
- Redis caching, async content preloading on startup, SQLite summary storage
- i18n support (zh-CN / en) via `LOCALE` env variable
- Dark mode and mobile responsive
- `sitemap.xml` / `robots.txt`

## Quick Start

### Docker Compose (Recommended)

```bash
git clone https://github.com/harveyTon/trilium-blog.git
cd trilium-blog

# Create config file
cp .env.example .env
# Edit .env, fill in TRILIUM_API_URL, TRILIUM_TOKEN, etc.
vim .env

docker compose up -d --build
```

Visit `http://localhost:8080` (port can be changed via `PORT` in `.env`).

The default deployment uses Docker Compose. The `docker-compose.yml` starts:

- `trilium-blog`
- `redis`

Notes:

- Redis address is fixed to `redis:6379` inside the container, no extra configuration needed.
- AI / code summary SQLite database is located at `/app/data/summaries.db` inside the container.
- The `trilium-blog-data` volume persists the summary database, preventing full AI summary regeneration after `docker compose up -d --build`.

### Local Development

**Backend** (requires Go 1.25+):

```bash
cd backend
go mod download
# Set environment variables or export TRILIUM_API_URL=... TRILIUM_TOKEN=... etc.
go run main.go
```

**Frontend** (requires Node 24+):

```bash
cd frontend
npm install
npm run dev
```

## Configuration

All configuration is managed via environment variables (`.env` file):

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `TRILIUM_API_URL` | Yes | — | Trilium ETAPI address |
| `TRILIUM_TOKEN` | Yes | — | ETAPI Token |
| `BLOG_TITLE` | No | — | Blog title |
| `BLOG_SUBTITLE` | No | — | Blog subtitle |
| `DOMAIN` | No | — | Blog domain, used for sitemap, page links, and internal resource detection |
| `ARTICLES_PER_PAGE` | No | `9` | Articles per page |
| `PORT` | No | `8080` | Server port |
| `LOCALE` | No | `zh-CN` | Blog language, supports `zh-CN` (Chinese) and `en` (English) |
| `IMAGE_PROXY_ENABLED` | No | `false` | Enable external image proxy |
| `IMAGE_PROXY_BASE_URL` | No | — | External image proxy URL (leave empty to use built-in `/api/imageproxy`) |
| `AI_SUMMARY_ENABLED` | No | `false` | Enable summary subsystem |
| `AI_SUMMARY_PROVIDER` | No | `openai-compatible` | AI provider type |
| `AI_SUMMARY_BASE_URL` | No | — | AI API base URL |
| `AI_SUMMARY_API_KEY` | No | — | AI API key |
| `AI_SUMMARY_MODEL` | No | — | AI model name |
| `AI_SUMMARY_PROMPT` | No | Built-in default | AI summary system prompt |
| `AI_SUMMARY_MODE` | No | `code` | `code` generates local summary only, `ai` keeps code summary and async generates AI summary |
| `AI_SUMMARY_CONCURRENCY` | No | `2` | AI summary concurrent workers |
| `AI_SUMMARY_RATE_LIMIT_MS` | No | `1200` | AI summary request interval (ms) |
| `AI_SUMMARY_TIMEOUT_MS` | No | `60000` | Single AI request timeout (ms) |
| `AI_SUMMARY_MAX_INPUT_CHARS` | No | `12000` | Max characters sent to AI |

The AI summary SQLite file is managed internally by the backend; in Docker it defaults to `/app/data/summaries.db` and is persisted via the `docker-compose.yml` volume.

## AI Summary

The project includes two layers of summary capability:

- `code summary`: Synchronously generated and stored, ensuring list and article pages always have a displayable summary.
- `ai summary`: Asynchronously generated and backfilled to the database, never blocking article content loading.

AI behavior:

- Article detail pages load content first.
- If AI summary is enabled and the article has no AI result yet, the backend queues a generation task immediately.
- The frontend shows a "generating" state for the AI summary card and polls for the result.
- Once the AI summary is ready, the article page updates automatically; list pages and featured posts using AI summary show an `AI` badge.
- AI generation context includes both article title and content.
- The `posts` API returns `summaries` directly; the frontend reuses existing results without re-requesting the summary endpoint.

To keep only local summaries without AI requests:

```env
AI_SUMMARY_ENABLED=false
```

or

```env
AI_SUMMARY_MODE=code
```

## Caching & Preloading

- All article lists, content, and attachments are cached via Redis (default TTL 5 minutes).
- On startup, all `#blog=true` article content is preloaded into Redis asynchronously; first visits hit cache directly without waiting for Trilium ETAPI.
- Preloading only caches raw content and does not trigger code summary or AI summary generation.
- If Redis is unavailable, the system falls back to no-cache mode, forwarding all requests directly to Trilium.

## Homepage & Article Pages

### Homepage

- The homepage consists of "Featured Posts" and "Latest Posts" sections.
- Featured posts come from the `#blogtop=true` label.
- Latest posts support pagination, which updates the browser URL.
- Page 1 title is `BLOG_TITLE | BLOG_SUBTITLE`.
- Page 2 and beyond append the page number, e.g., `BLOG_TITLE | BLOG_SUBTITLE - Page 2`.

### Article Page

- Article content, TOC, summary, and code block info are returned by the backend.
- Article page title format: `Article Title - BLOG_TITLE | BLOG_SUBTITLE`.
- Normal mode shows summary, content, and source link.
- Reading mode shares the same route, switching via page-level class.
- Reading mode preferences are persisted, including:
  - TOC collapsed state
  - Width
  - Density
  - Font size
  - Reading theme

## Code Blocks & Highlighting

- The backend extracts code block metadata when processing article HTML.
- Code language detection priority:
  1. Trilium MIME class name (built-in 120+ MIME type mapping table)
  2. [Chroma](https://github.com/alecthomas/chroma) lexical analysis
  3. [enry](https://github.com/go-enry/go-enry) statistical classifier
  4. Fallback to `plaintext`
- Frontend uses [Shiki](https://github.com/shikijs/shiki) for syntax highlighting.
- Supports code language labels, copy button, line numbers, and light/dark theme switching.

## Usage

In Trilium Notes:

- Add `#blog=true` to notes you want to publish
- Add `#blogtop=true` to notes you want to feature

The blog will automatically read and display these notes.

## Tech Stack

**Backend:** Go 1.25 / Gin / goquery / bluemonday / Redis / SQLite / Chroma / enry

**Frontend:** Vue 3 / Vite / Element Plus / Pinia / Vue Router / Fancybox / Shiki

## Acknowledgments

- [Trilium](https://github.com/TriliumNext/Trilium)
- [Trilium Notes](https://github.com/zadam/trilium)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Vue.js](https://vuejs.org/)
- [Element Plus](https://element-plus.org/)
- [Fancybox](https://fancyapps.com/fancybox/)
- [Shiki](https://github.com/shikijs/shiki)
- [Chroma](https://github.com/alecthomas/chroma)
- [enry](https://github.com/go-enry/go-enry)

## License

[MIT](LICENSE)
