package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harveyTon/trilium-blog/backend/config"
	"github.com/harveyTon/trilium-blog/backend/pkg/logger"
	"github.com/harveyTon/trilium-blog/backend/services"
)

// GetArticles godoc
// @Summary Get a list of articles
// @Description Get a paginated list of blog articles
// @Tags articles
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Success 200 {object} map[string]interface{}
// @Router /articles [get]
func GetArticles(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	articlesPerPage := config.Config.ArticlesPerPage
	logger.Infof("Fetching articles for page %d with %d articles per page", page, articlesPerPage)

	articles, totalArticles, err := services.GetArticles(page, articlesPerPage)
	if err != nil {
		logger.Infof("Error fetching articles: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles"})
		return
	}

	totalPages := (totalArticles + articlesPerPage - 1) / articlesPerPage

	logger.Infof("Fetched %d articles, total articles: %d, total pages: %d", len(articles), totalArticles, totalPages)

	c.JSON(http.StatusOK, gin.H{
		"articles":        articles,
		"currentPage":     page,
		"totalPages":      totalPages,
		"totalArticles":   totalArticles,
		"articlesPerPage": articlesPerPage,
	})
}

// GetArticle godoc
// @Summary Get a single article
// @Description Get a single blog article by its noteId
// @Tags articles
// @Accept json
// @Produce json
// @Param noteId path string true "Note ID"
// @Success 200 {object} models.Article
// @Router /articles/{noteId} [get]
func GetArticle(c *gin.Context) {
	noteId := c.Param("noteId")

	article, content, err := services.GetArticle(noteId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch article"})
		return
	}

	var pageUrl string
	for _, attr := range article.Attributes {
		if attr.Name == "pageUrl" {
			pageUrl = attr.Value
			break
		}
	}

	data := gin.H{
		"title":        article.Title,
		"blogName":     config.Config.BlogName,
		"domain":       config.Config.Domain,
		"noteId":       noteId,
		"content":      template.HTML(content),
		"pageUrl":      pageUrl,
		"dateModified": article.DateModified,
	}

	c.JSON(http.StatusOK, data)
}

func GenerateSitemap(c *gin.Context) {
	sitemap, err := services.GenerateSitemap()
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to generate sitemap")
		return
	}

	c.Header("Content-Type", "application/xml")
	c.String(http.StatusOK, sitemap)
}

func RobotsTxt(c *gin.Context) {
	c.String(http.StatusOK, `User-agent: *
Allow: /
Sitemap: /sitemap.xml`)
}
