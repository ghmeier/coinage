package handlers

import (
	"gopkg.in/gin-gonic/gin.v1"
)

func empty() *gin.H {
	return &gin.H{"success": true}
}

func errResponse(message string) *gin.H {
	return &gin.H{"error": message}
}
