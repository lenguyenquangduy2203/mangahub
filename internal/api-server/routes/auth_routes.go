package routes

import (
	"mangahub/internal/api-server/handlers"

	"github.com/gin-gonic/gin"
)

func AuthRoutesV1(router *gin.RouterGroup, handler *handlers.AuthHandler) {
	// Register
	router.POST("/register", handler.RegisterV1)

	// Login
	router.POST("/login", handler.LoginV1)
}
