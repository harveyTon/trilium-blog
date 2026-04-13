package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/harveyTon/trilium-blog/backend/blog"
	"github.com/harveyTon/trilium-blog/backend/config"
	"github.com/harveyTon/trilium-blog/backend/etapi"
	"github.com/harveyTon/trilium-blog/backend/handlers"
	"github.com/harveyTon/trilium-blog/backend/pkg/logger"
)

func setupRouter(apiHandler *handlers.APIHandler) *gin.Engine {
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
		api.GET("/assets/:attachmentId", apiHandler.GetAsset)
		api.GET("/imageproxy", apiHandler.ImageProxy)
	}
	r.GET("/sitemap.xml", apiHandler.Sitemap)
	r.GET("/robots.txt", apiHandler.Robots)

	staticDir := resolveFrontendDist()
	r.Static("/assets", filepath.Join(staticDir, "assets"))
	r.StaticFile("/favicon.ico", filepath.Join(staticDir, "favicon.ico"))
	r.StaticFile("/logo.png", filepath.Join(staticDir, "logo.png"))

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

func main() {
	logger.Init()
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	config.LoadConfig()

	etapiClient := etapi.NewClient(config.Config.TriliumApiUrl, config.Config.TriliumToken)
	summaryStore, err := blog.NewSummaryStoreDB(config.Config.AISummary.DatabasePath)
	if err != nil {
		logger.Fatal("Failed to initialize summary store", err)
	}
	defer summaryStore.Close()

	var aiQueue *blog.AISummaryQueue
	if config.Config.AISummary.Enabled {
		aiQueue = blog.NewAISummaryQueue(
			summaryStore,
			config.Config.AISummary.BaseURL,
			config.Config.AISummary.APIKey,
			config.Config.AISummary.Model,
			config.Config.AISummary.Prompt,
			config.Config.AISummary.Concurrency,
			config.Config.AISummary.RateLimitMs,
		)
	}

	service := blog.NewService(
		etapiClient,
		&blog.NoopStore{},
		blog.WithBlogName(config.Config.BlogName),
		blog.WithBlogTitle(config.Config.BlogTitle),
		blog.WithDomain(config.Config.Domain),
		blog.WithPageSize(config.Config.ArticlesPerPage),
		blog.WithImageProxyEnabled(config.Config.ImageProxy.Enabled),
		blog.WithImageProxyBaseUrl(config.Config.ImageProxy.BaseURL),
		blog.WithSummaryStore(summaryStore),
		blog.WithAISummaryQueue(aiQueue),
		blog.WithAISummaryEnabled(config.Config.AISummary.Enabled),
	)
	apiHandler := handlers.NewAPIHandler(service)
	r := setupRouter(apiHandler)

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
