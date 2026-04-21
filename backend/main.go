package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harveyTon/trilium-blog/backend/blog"
	"github.com/harveyTon/trilium-blog/backend/config"
	"github.com/harveyTon/trilium-blog/backend/etapi"
	"github.com/harveyTon/trilium-blog/backend/handlers"
	"github.com/harveyTon/trilium-blog/backend/pkg/logger"
)

const (
	defaultRedisAddr       = "redis:6379"
	defaultRedisPassword   = ""
	defaultRedisDB         = 0
	defaultRedisTTLSeconds = 300
	customAssetsDir        = "./custom"
)

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func setupRouter(apiHandler *handlers.APIHandler, staticDir string) *gin.Engine {
	gin.DisableConsoleColor()
	r := gin.Default()
	r.Use(logger.GinLogger())
	r.Use(gin.Recovery())

	api := r.Group("/api")
	{
		api.GET("/site", apiHandler.GetSite)
		api.GET("/posts", apiHandler.ListPosts)
		api.GET("/posts/featured", apiHandler.ListFeaturedPosts)
		api.GET("/search", apiHandler.SearchPosts)
		api.GET("/posts/:noteId", apiHandler.GetPost)
		api.GET("/posts/:noteId/summary", apiHandler.GetPostSummary)
		api.GET("/assets/:attachmentId", apiHandler.GetAsset)
		api.GET("/imageproxy", apiHandler.ImageProxy)
		api.GET("/health", healthCheck)
	}
	admin := r.Group("/api/admin")
	admin.Use(apiHandler.AdminAuthMiddleware)
	{
		admin.GET("/cache/stats", apiHandler.CacheStats)
		admin.POST("/cache/invalidate", apiHandler.InvalidateCache)
		admin.POST("/cache/preload", apiHandler.TriggerPreload)
	}
	if config.Config.AdminToken != "" {
		r.GET("/admin", apiHandler.AdminPage)
	}
	r.GET("/sitemap.xml", apiHandler.Sitemap)
	r.GET("/robots.txt", apiHandler.Robots)

	r.Static("/assets", filepath.Join(staticDir, "assets"))
	r.StaticFile("/favicon.ico", resolveStaticFile(staticDir, "favicon.ico"))
	r.StaticFile("/logo.png", resolveStaticFile(staticDir, "logo.png"))

	r.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/api") && c.Request.URL.Path != "/sitemap.xml" && c.Request.URL.Path != "/robots.txt" {
			c.File(filepath.Join(staticDir, "index.html"))
		} else {
			c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
		}
	})

	return r
}

func resolveFrontendDist() string {
	paths := []string{
		"./frontend/dist",
		"../frontend/dist",
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return "./frontend/dist"
}

func dataDir() string {
	if d := os.Getenv("DATA_DIR"); d != "" {
		return d
	}
	return "./data"
}

func resolveStaticFile(staticDir, name string) string {
	customPath := filepath.Join(customAssetsDir, name)
	if _, err := os.Stat(customPath); err == nil {
		logger.Info(fmt.Sprintf("Using custom %s from %s", name, customPath))
		return customPath
	}
	return filepath.Join(staticDir, name)
}

func runStartupChecks(etapiClient *etapi.Client, dir, staticDir string) {
	logger.Info("========== Startup Checks ==========")

	logger.Info(fmt.Sprintf("[Config] TRILIUM_API_URL = %s", config.Config.TriliumApiUrl))
	logger.Info(fmt.Sprintf("[Config] BLOG_TITLE = %q", config.Config.BlogTitle))
	logger.Info(fmt.Sprintf("[Config] BLOG_SUBTITLE = %q", config.Config.BlogSubtitle))
	logger.Info(fmt.Sprintf("[Config] DOMAIN = %s", config.Config.Domain))
	logger.Info(fmt.Sprintf("[Config] LOCALE = %s", config.Config.Locale))
	logger.Info(fmt.Sprintf("[Config] ARTICLES_PER_PAGE = %d", config.Config.ArticlesPerPage))
	logger.Info(fmt.Sprintf("[Config] ADMIN_TOKEN = %s", boolStr(config.Config.AdminToken != "")))
	logger.Info(fmt.Sprintf("[Config] IMAGE_PROXY = enabled=%v, base_url=%s", config.Config.ImageProxy.Enabled, config.Config.ImageProxy.BaseURL))
	logger.Info(fmt.Sprintf("[Config] AI_SUMMARY = enabled=%v, mode=%s, provider=%s", config.Config.AISummary.Enabled, config.Config.AISummary.Mode, config.Config.AISummary.Provider))

	logger.Info(fmt.Sprintf("[Data] Data directory: %s", dir))
	if info, err := os.Stat(dir); err != nil {
		logger.Error(fmt.Sprintf("[Data] Data directory not accessible: %s", dir), err)
	} else {
		logger.Info(fmt.Sprintf("[Data] Data directory OK (mode=%s)", info.Mode().Perm()))
	}

	logger.Info(fmt.Sprintf("[Frontend] Static directory: %s", staticDir))
	if _, err := os.Stat(staticDir); err != nil {
		logger.Error(fmt.Sprintf("[Frontend] Static directory not found: %s", staticDir), err)
	} else {
		indexPath := filepath.Join(staticDir, "index.html")
		if _, err := os.Stat(indexPath); err != nil {
			logger.Error("[Frontend] index.html not found in static directory", err)
		} else {
			logger.Info("[Frontend] index.html OK")
		}
	}

	customAssets := []string{"favicon.ico", "logo.png"}
	for _, name := range customAssets {
		customPath := filepath.Join(customAssetsDir, name)
		if _, err := os.Stat(customPath); err == nil {
			logger.Info(fmt.Sprintf("[Assets] Custom %s detected: %s", name, customPath))
		}
	}

	logger.Info("[Trilium] Testing connectivity...")
	start := time.Now()
	notes, err := etapiClient.GetNotes("#blog=true")
	elapsed := time.Since(start).Round(time.Millisecond)
	if err != nil {
		logger.Error(fmt.Sprintf("[Trilium] Connection failed (%s): %v", elapsed, err), err)
		logger.Error("[Trilium] Please verify TRILIUM_API_URL and TRILIUM_TOKEN are correct and Trilium is running", err)
	} else {
		logger.Info(fmt.Sprintf("[Trilium] Connection OK (%s, %d blog posts found)", elapsed, len(notes)))
		if len(notes) == 0 {
			logger.Warn("[Trilium] No notes with #blog=true label found. Make sure notes are tagged in Trilium.")
		}
	}

	logger.Info(fmt.Sprintf("[AI Summary] provider=%s base_url=%s model=%s",
		config.Config.AISummary.Provider,
		maskSecret(config.Config.AISummary.BaseURL),
		config.Config.AISummary.Model))
	if config.Config.AISummary.AIRequestsEnabled() {
		if config.Config.AISummary.BaseURL == "" {
			logger.Warn("[AI Summary] Mode is 'ai' but AI_SUMMARY_BASE_URL is empty")
		}
		if config.Config.AISummary.APIKey == "" {
			logger.Warn("[AI Summary] Mode is 'ai' but AI_SUMMARY_API_KEY is empty")
		}
		if config.Config.AISummary.Model == "" {
			logger.Warn("[AI Summary] Mode is 'ai' but AI_SUMMARY_MODEL is empty")
		}
	}

	logger.Info("========== Startup Checks Complete ==========")
}

func boolStr(v bool) string {
	if v {
		return "set"
	}
	return "not set"
}

func maskSecret(s string) string {
	if s == "" {
		return "(empty)"
	}
	return s
}

func main() {
	config.LoadConfig()
	logger.Init(config.Config.LogLevel)
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	dir := dataDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Error("Failed to create data directory", err)
	}

	etapiClient := etapi.NewClient(config.Config.TriliumApiUrl, config.Config.TriliumToken)
	var err error

	cacheStore := blog.Store(&blog.NoopStore{})
	var redisStore *blog.RedisStore
	redisStore, err = blog.NewRedisStore(
		defaultRedisAddr,
		defaultRedisPassword,
		defaultRedisDB,
		defaultRedisTTLSeconds,
	)
	if err != nil {
		logger.Warn("Redis unavailable; falling back to file cache")
		cacheStore = blog.InitFileCache(dir)
	} else {
		cacheStore = redisStore
		defer redisStore.Close()
	}

	summaryDatabasePath := filepath.Join(dir, "summaries.db")
	summaryStore, err := blog.NewSummaryStoreDB(summaryDatabasePath)
	if err != nil {
		logger.Error("Failed to initialize summary store; continuing without persisted summaries", err)
	}
	if summaryStore != nil {
		defer summaryStore.Close()
	}

	var aiQueue *blog.AISummaryQueue
	aiSummaryEnabled := summaryStore != nil && config.Config.AISummary.AIRequestsEnabled()
	if aiSummaryEnabled {
		aiQueue = blog.NewAISummaryQueue(
			summaryStore,
			config.Config.AISummary.Provider,
			config.Config.AISummary.BaseURL,
			config.Config.AISummary.APIKey,
			config.Config.AISummary.Model,
			config.Config.AISummary.Prompt,
			config.Config.AISummary.Concurrency,
			config.Config.AISummary.RateLimitMs,
			config.Config.AISummary.TimeoutMs,
			config.Config.AISummary.MaxInputChars,
		)
	}

	service := blog.NewService(
		etapiClient,
		cacheStore,
		blog.WithBlogTitle(config.Config.BlogTitle),
		blog.WithBlogSubtitle(config.Config.BlogSubtitle),
		blog.WithDomain(config.Config.Domain),
		blog.WithLocale(config.Config.Locale),
		blog.WithPageSize(config.Config.ArticlesPerPage),
		blog.WithImageProxyEnabled(config.Config.ImageProxy.Enabled),
		blog.WithImageProxyBaseUrl(config.Config.ImageProxy.BaseURL),
		blog.WithSummaryStore(summaryStore),
		blog.WithAISummaryQueue(aiQueue),
		blog.WithAISummaryEnabled(aiSummaryEnabled),
	)

	staticDir := resolveFrontendDist()
	runStartupChecks(etapiClient, dir, staticDir)

	apiHandler := handlers.NewAPIHandler(service, config.Config.AdminToken, config.Config.Locale)
	r := setupRouter(apiHandler, staticDir)

	go service.Preload()

	defaultPort := os.Getenv("PORT")
	if defaultPort == "" {
		defaultPort = "8080"
	}
	if len(os.Args) > 1 {
		defaultPort = os.Args[1]
	}
	logger.Info(fmt.Sprintf("Server started on port %s", defaultPort))
	logger.Fatal("Server failed to start", r.Run(":"+defaultPort))
}
