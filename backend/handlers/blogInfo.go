package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harveyTon/trilium-blog/backend/config"
)

// GetBlogInfo godoc
// @Summary Get blog information
// @Description Retrieve the blog's name and title from the configuration
// @Tags blog
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Blog information"
// @Router /blog/info [get]
func GetBlogInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"blogName":  config.Config.BlogName,
		"blogTitle": config.Config.BlogTitle,
	})
}
