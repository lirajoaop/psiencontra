package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func sendError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

func sendSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}
