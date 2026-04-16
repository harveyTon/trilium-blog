# AGENTS.md

## Commands

```bash
# Backend tests (run from backend/)
go test ./...

# Backend single package
go test ./blog/...
go test ./handlers/...

# Frontend build (run from frontend/)
npm run build

# Frontend dev server
npm run dev

# Full rebuild + restart (run from repo root)
docker compose up -d --build
```

There is no lint or typecheck command configured. Run `go test` and `npm run build` as verification before committing.

## Architecture

Monorepo with two independent packages:

- **`backend/`** — Go (Gin), serves API + embedded frontend SPA
- **`frontend/`** — Vue 3 (Vite), builds to `frontend/dist/`

The Go binary embeds the frontend dist at runtime via `resolveFrontendDist()` (checks `./frontend/dist` then `../frontend/dist`). In Docker, the built frontend is copied into the Go container at `./frontend/dist`.

### Backend structure

```
backend/
  main.go              — entrypoint, Gin router, config loading, async preload
  config/config.go     — AppConfig, reads environment variables (including LOCALE)
  etapi/client.go      — Trilium ETAPI HTTP client
  blog/service.go      — core logic: post listing, featured posts, content processing, TOC extraction, preload, code language detection
  blog/search.go       — search matching, scoring, snippet shaping
  blog/summary*.go     — summary persistence, queue, AI summary flow
  blog/models.go       — Post, search, summary, site structs (JSON API contract)
  blog/cache.go        — NoopStore / RedisStore / FileStore implementations
  handlers/api.go      — Gin handlers for /api/* routes
```

Request flow: `handlers` → `blog.Service` → `etapi.Client` → Trilium Notes ETAPI.

### Frontend structure

```
frontend/src/
  App.vue       — app shell, CSS variables, keep-alive wiring
  i18n/         — lightweight i18n: index.js (t/setLocale/locale), locales/zh-CN.js + en.js
  router/       — routes: / → HomePage, /search → SearchPage, /post/:noteId → Article, * → NotFound
  store/        — Pinia stores for site/search state
  api/          — blog/search/summary API modules
  components/   — app/home/search/article component layers
  composables/  — search, summary, reading, article enhancement logic
  views/
    Home.vue    — search entry + featured section + paginated post feed
    Search.vue  — dedicated search page
    Article.vue — article detail shell, reading progress, TOC, reading mode settings
```

### API routes

| Route | Purpose |
|---|---|
| `GET /api/site` | Blog config (name, domain, locale, imageProxy) |
| `GET /api/posts?page=N` | Paginated post list (default 9/page) |
| `GET /api/posts/featured` | Featured post list from `blogtop=true` |
| `GET /api/search?q=...` | Site search and preview results |
| `GET /api/posts/:noteId` | Single post with contentHtml + TOC |
| `GET /api/posts/:noteId/summary` | Post summaries |
| `GET /api/assets/:attachmentId` | Proxy image from Trilium |
| `GET /api/imageproxy?url=` | Internal fallback proxy for external images |

All non-API, non-static routes serve `index.html` (SPA fallback).

## Key gotchas

- **Attachment URLs**: Trilium outputs relative paths like `api/attachments/{id}/image/{file}`. `blog/service.go:extractAttachmentId()` rewrites these to `/api/assets/{id}` using regex — the ID is only the first path segment.
- **ETAPI search**: The `#blog=true` label search parameter must be URL-encoded (the `#` would be stripped as a fragment). Handled in `etapi/client.go`.
- **TOC heading IDs**: `extractTOC()` returns both the TOC array and the modified HTML with `id` attributes injected. The modified HTML must flow into `processContent()` — do not pass the original HTML.
- **CSS variables**: All colors defined in `App.vue :root` / `html.dark`. No hardcoded colors elsewhere.
- **Language detection chain**: Trilium MIME class mapping table (`triliumClassToLangID` in service.go, 120+ entries) → Chroma Analyse → enry classifier → plaintext fallback. `language-auto`, `language-text-x-trilium-auto`, `language-text`, `language-plain` are ignored (treated as no class) to trigger content-based detection.
- **i18n**: No vue-i18n dependency. Custom lightweight module at `frontend/src/i18n/index.js` exports `t(key, params)`, `setLocale(locale)`, `locale()`. Locale comes from backend `/api/site` response, set once on site load via Pinia store. Translation files are plain JS objects in `frontend/src/i18n/locales/`.
- **Options API + i18n**: Components using `setup()` must return `{ t }` from setup. Components without `setup()` must include `methods: { t }` for template access. This is the single most common bug when adding i18n to a new component.
- **Theme**: Dark mode toggled via `useDark()` (VueUse), persisted in `vueuse-color-scheme` localStorage.
- **Keep-alive**: `HomePage` is wrapped in `<keep-alive>` to preserve pagination state on back navigation.
- **Preload**: On startup, `service.Preload()` runs async to cache all post content. Uses Redis by default; falls back to file-based cache (`DATA_DIR/cache/`) if Redis is unavailable. It does not trigger summary generation.
- **AI summary config**: Summary configuration is env-driven. `AI_SUMMARY_ENABLED=true` only opens the subsystem; actual async AI requests run only when `AI_SUMMARY_MODE=ai`.

## Config

- All runtime config is env-driven. Local development and Docker both use `.env` / exported environment variables.
- Copy `.env.example` to `.env` before local or Docker startup.
- Blog posts are identified by `#blog=true` label in Trilium Notes. Featured posts use `#blogtop=true`.
