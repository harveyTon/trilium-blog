package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/harveyTon/trilium-blog/backend/config"
	"github.com/harveyTon/trilium-blog/backend/handlers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/harveyTon/trilium-blog/backend/docs"
	"github.com/harveyTon/trilium-blog/backend/pkg/logger"
)

// @title Trilium Blog API
// @version 1.0
// @description This is a blog API server for Trilium Notes.
// @host localhost:8080
// @BasePath /api
func setupRouter() *gin.Engine {
	gin.DisableConsoleColor()
	r := gin.Default()
	r.Use(logger.GinLogger())
	r.Use(gin.Recovery())

	api := r.Group("/api")
	{
		api.GET("/info", handlers.GetBlogInfo)
		api.GET("/articles", handlers.GetArticles)
		api.GET("/articles/:noteId", handlers.GetArticle)
		api.GET("/attachments/:attachmentId", handlers.GetAttachment)

	}
	r.GET("/attachments/:attachmentId", handlers.GetAttachment)
	r.GET("/sitemap.xml", handlers.GenerateSitemap)
	r.GET("/robots.txt", handlers.RobotsTxt)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	_ = mime.AddExtensionType(".js", "application/javascript")
	_ = mime.AddExtensionType(".mjs", "application/javascript")

	r.Static("/assets", "./frontend/dist/assets")
	r.StaticFile("/favicon.ico", "./frontend/dist/favicon.ico")
	r.StaticFile("/logo.png", "./frontend/dist/logo.png")

	r.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.File("./frontend/dist/index.html")
		} else {
			c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
		}
	})
	r.Use(func(c *gin.Context) {
		if strings.HasSuffix(c.Request.URL.Path, ".js") || strings.HasSuffix(c.Request.URL.Path, ".mjs") {
			c.Writer.Header().Set("Content-Type", "application/javascript")
		}
		c.Next()
	})

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	logger.Init()
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	config.LoadConfig()

	r := setupRouter()

	defaultPort := "8080"
	if len(os.Args) > 1 {
		defaultPort = os.Args[1]
	}
	logger.Info(fmt.Sprintf("Server started on port %s", defaultPort))
	logger.Fatal("Server failed to start", r.Run(":"+defaultPort))

}
