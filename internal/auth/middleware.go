package auth

import (
	"mangahub/pkg/models/dtos"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var USER_KEY = "auth_user_token"

//func AuthMiddleware(tokenManager TokenManager) gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		authHeader := ctx.GetHeader("Authorization")
//
//		if authHeader == "" {
//			ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
//				Code:    "INVALID_TOKEN",
//				Message: "Authorization header required",
//			})
//			ctx.Abort()
//			return
//		}
//
//		parts := strings.Split(authHeader, " ")
//		if len(parts) != 2 || parts[0] != "Bearer" {
//			ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
//				Code:    "INVALID_TOKEN",
//				Message: "Invalid Authorization header format",
//			})
//			ctx.Abort()
//			return
//		}
//
//		tokenString := parts[1]
//
//		claims, err := tokenManager.ValidateToken(tokenString)
//		if err != nil {
//			ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
//				Code:    "INVALID_TOKEN",
//				Message: "Invalid or expired token",
//			})
//			ctx.Abort()
//			return
//		}
//
//		// Store the user claims in the context for downstream handlers
//		ctx.Set(USER_KEY, claims)
//
//		ctx.Next()
//	}
//}

func AuthMiddleware(tokenManager TokenManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var tokenString string

		authHeader := ctx.GetHeader("Authorization")

		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
					Code:    "INVALID_TOKEN",
					Message: "Invalid Authorization header format",
				})
				ctx.Abort()
				return
			}
			tokenString = parts[1]
		} else {
			tokenString = ctx.Query("token")
		}

		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
				Code:    "INVALID_TOKEN",
				Message: "Authorization header or token query parameter required",
			})
			ctx.Abort()
			return
		}

		claims, err := tokenManager.ValidateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
				Code:    "INVALID_TOKEN",
				Message: "Invalid or expired token",
			})
			ctx.Abort()
			return
		}

		ctx.Set(USER_KEY, claims)
		ctx.Next()
	}
}
