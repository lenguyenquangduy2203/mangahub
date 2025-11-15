package routes

import (
	"mangahub/internal/api-server/handlers"

	"github.com/gin-gonic/gin"
)

func UserRoutesV1(router *gin.RouterGroup, handler *handlers.UserHandler) {
	// Add manga to user's library
	router.POST("/library", handler.AddMangaV1)
	// Get user's library
	router.GET("/library", handler.GetUserLibraryV1)
	// Update user's progress
	router.PUT("/library", handler.UpdateReadingProgressV1)
}
