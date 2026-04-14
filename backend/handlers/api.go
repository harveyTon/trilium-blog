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

var blockedIPNets = mustParseCIDRs([]string{
	"100.64.0.0/10",
})

type APIHandler struct {
	service    *blog.Service
	adminToken string
}

func NewAPIHandler(service *blog.Service, adminToken string) *APIHandler {
	return &APIHandler{service: service, adminToken: adminToken}
}

func (h *APIHandler) AdminAuthMiddleware(c *gin.Context) {
	if h.adminToken == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin token not configured"})
		c.Abort()
		return
	}
	token := c.GetHeader("Authorization")
	if token == "" {
		token = c.Query("token")
	}
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}
	if token != h.adminToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid admin token"})
		c.Abort()
		return
	}
	c.Next()
}

type invalidateRequest struct {
	Scope string `json:"scope"`
	Type  string `json:"type"`
	ID    string `json:"id"`
}

func (h *APIHandler) InvalidateCache(c *gin.Context) {
	var req invalidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		req.Scope = "all"
	}

	switch req.Scope {
	case "note":
		if req.ID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required for note scope"})
			return
		}
		h.service.InvalidateNote(req.ID)
		c.JSON(http.StatusOK, gin.H{"invalidated": "note", "id": req.ID})
	case "notes-list":
		if req.ID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id (search key) is required for notes-list scope"})
			return
		}
		h.service.InvalidateNotesList(req.ID)
		c.JSON(http.StatusOK, gin.H{"invalidated": "notes-list", "id": req.ID})
	case "attachment":
		if req.ID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id is required for attachment scope"})
			return
		}
		h.service.InvalidateAttachment(req.ID)
		c.JSON(http.StatusOK, gin.H{"invalidated": "attachment", "id": req.ID})
	case "type":
		count := h.service.InvalidateByType(req.Type)
		c.JSON(http.StatusOK, gin.H{"invalidated": req.Type, "keys_removed": count})
	case "all", "":
		count := h.service.InvalidateAll()
		c.JSON(http.StatusOK, gin.H{"invalidated": "all", "keys_removed": count})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scope, use: all, type, note, notes-list, attachment"})
	}
}

func (h *APIHandler) CacheStats(c *gin.Context) {
	stats := h.service.GetCacheStats()
	c.JSON(http.StatusOK, stats)
}

func (h *APIHandler) TriggerPreload(c *gin.Context) {
	started := h.service.TriggerPreload()
	if !started {
		c.JSON(http.StatusConflict, gin.H{"error": "preload already in progress"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "preload started"})
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

func (h *APIHandler) GetPostSummary(c *gin.Context) {
	noteId := c.Param("noteId")

	summaries, err := h.service.GetPostSummaries(noteId)
	if err != nil {
		if _, ok := err.(*blog.BlogError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post summary"})
		return
	}

	c.JSON(http.StatusOK, summaries)
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

	host := strings.ToLower(parsedURL.Hostname())
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
				redirectHost := strings.ToLower(req.URL.Hostname())
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
	host = strings.TrimSpace(strings.ToLower(host))
	if host == "" {
		return true
	}
	if strings.Contains(host, ":") && !strings.Contains(host, "]") {
		if parsedHost, _, err := net.SplitHostPort(host); err == nil {
			host = parsedHost
		}
	}
	host = strings.TrimPrefix(host, "[")
	host = strings.TrimSuffix(host, "]")

	if host == "localhost" || strings.HasSuffix(host, ".local") {
		return true
	}
	if ip := net.ParseIP(host); ip != nil {
		return isPrivateIP(ip)
	}
	ips, err := net.LookupIP(host)
	if err != nil {
		return false
	}
	for _, ip := range ips {
		if isPrivateIP(ip) {
			return true
		}
	}
	return false
}

func isPrivateIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() || ip.IsMulticast() {
		return true
	}
	for _, block := range blockedIPNets {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func mustParseCIDRs(values []string) []*net.IPNet {
	nets := make([]*net.IPNet, 0, len(values))
	for _, value := range values {
		_, network, err := net.ParseCIDR(value)
		if err != nil {
			panic(err)
		}
		nets = append(nets, network)
	}
	return nets
}
