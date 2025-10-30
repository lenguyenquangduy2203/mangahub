package main

import (
	"net/http"
	"os"

	"mangahub/pkg/utils/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
}

func main() {
	router := gin.Default()

	// version 1
	apiV1 := router.Group("/v1")

	apiV1.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "good",
		})
	})

	// Read port from env
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "3000" // fallback default
	}

	router.Run(":" + port)
}
