package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/harveyTon/trilium-blog/backend/blog"
)

type APIHandler struct {
	service *blog.Service
}

func NewAPIHandler(service *blog.Service) *APIHandler {
	return &APIHandler{service: service}
}

func (h *APIHandler) GetSite(c *gin.Context) {
	site := h.service.GetSite()
	c.JSON(http.StatusOK, site)
}

func (h *APIHandler) ListPosts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	posts, err := h.service.ListPosts(page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *APIHandler) GetPost(c *gin.Context) {
	noteId := c.Param("noteId")

	post, err := h.service.GetPost(noteId)
	if err != nil {
		if _, ok := err.(*blog.BlogError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *APIHandler) GetAsset(c *gin.Context) {
	attachmentId := c.Param("attachmentId")

	content, contentType, err := h.service.GetAsset(attachmentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch asset"})
		return
	}

	c.Data(http.StatusOK, contentType, content)
}

func (h *APIHandler) Sitemap(c *gin.Context) {
	c.Header("Content-Type", "application/xml")
	c.String(http.StatusOK, `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"></urlset>`)
}

func (h *APIHandler) Robots(c *gin.Context) {
	c.String(http.StatusOK, "User-agent: *\nAllow: /\nSitemap: /sitemap.xml")
}
