package handlers

import (
	"strconv"

	"gopkg.in/gin-gonic/gin.v1"
)

func empty() *gin.H {
	return &gin.H{"success": true}
}

func getPaging(ctx *gin.Context) (int, int) {
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
	return offset, limit
}

func errResponse(message string) *gin.H {
	return &gin.H{"error": message}
}
