package routes

import (
	"mangahub/internal/api-server/handlers"

	"github.com/gin-gonic/gin"
)

func MangaRoutesV1(router *gin.RouterGroup, hander *handlers.MangaHandler) {
	// Query mangas
	router.GET("", hander.FindMangasV1)

	// Get manga details
	router.GET("/:id", hander.GetMangaV1)
}
