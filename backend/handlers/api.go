package handlers

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

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

func (h *APIHandler) SearchPosts(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	preview := c.DefaultQuery("preview", "false") == "true"
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 5
	}

	result, err := h.service.SearchPosts(query, preview, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search posts"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *APIHandler) ListFeaturedPosts(c *gin.Context) {
	posts, err := h.service.ListFeaturedPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch featured posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": posts})
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
	sitemap, err := h.service.GenerateSitemap()
	if err != nil {
		c.String(http.StatusInternalServerError, `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"></urlset>`)
		return
	}
	c.String(http.StatusOK, sitemap)
}

func (h *APIHandler) Robots(c *gin.Context) {
	c.String(http.StatusOK, "User-agent: *\nAllow: /\nSitemap: /sitemap.xml")
}

func (h *APIHandler) ImageProxy(c *gin.Context) {
	rawURL := c.Query("url")
	if rawURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		return
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil || parsedURL.Host == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only http and https are allowed"})
		return
	}

	host := strings.ToLower(parsedURL.Host)
	if isBlockedHost(host) {
		c.JSON(http.StatusForbidden, gin.H{"error": "proxying this host is not allowed"})
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return http.ErrUseLastResponse
			}
			if via != nil {
				redirectHost := strings.ToLower(req.URL.Host)
				if isBlockedHost(redirectHost) {
					return http.ErrUseLastResponse
				}
			}
			return nil
		},
	}

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid url"})
		return
	}
	req.Header.Set("User-Agent", "trilium-blog-image-proxy/1.0")

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to fetch image"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{"error": "upstream returned non-200 status"})
		return
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "response is not an image"})
		return
	}

	body, err := io.ReadAll(http.MaxBytesReader(nil, resp.Body, 10*1024*1024))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "response too large or read error"})
		return
	}

	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=86400")
	c.Data(http.StatusOK, contentType, body)
}

func isBlockedHost(host string) bool {
	if host == "localhost" || host == "127.0.0.1" || host == "::1" || strings.HasSuffix(host, ".local") {
		return true
	}
	if ip := net.ParseIP(host); ip != nil {
		return isPrivateIP(ip)
	}
	return false
}

func isPrivateIP(ip net.IP) bool {
	blocks := []string{
		"10.",
		"172.16.", "172.17.", "172.18.", "172.19.", "172.20.",
		"172.21.", "172.22.", "172.23.", "172.24.", "172.25.",
		"172.26.", "172.27.", "172.28.", "172.29.", "172.30.", "172.31.",
		"192.168.",
	}
	for _, block := range blocks {
		if strings.HasPrefix(ip.String(), block) {
			return true
		}
	}
	return false
}
