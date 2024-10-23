package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/harveyTon/trilium-blog/backend/services"
)

// GetAttachment godoc
// @Summary Get attachment content
// @Description Retrieve the content of the specified attachment by its ID, returning only attachments that belong to blog posts
// @Tags attachments
// @Accept json
// @Produce octet-stream
// @Param attachmentId path string true "Attachment ID"
// @Success 200 {array} byte "Attachment content"
// @Router /attachments/{attachmentId} [get]
func GetAttachment(c *gin.Context) {
	attachmentId := c.Param("attachmentId")

	content, contentType, err := services.GetAttachment(attachmentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attachment"})
		return
	}

	c.Data(http.StatusOK, contentType, content)
}
