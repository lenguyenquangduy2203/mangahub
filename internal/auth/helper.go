package auth

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func GetUserIDFromContext(ctx *gin.Context) (string, error) {
	value, exists := ctx.Get(USER_KEY)
	if !exists {
		return "", errors.New("user claims not found in context. Middleware not run?")
	}

	claims, ok := value.(*JWTClaims)
	if !ok {
		return "", errors.New("context value is not of expected *JWTClaims type")
	}

	return claims.UserID, nil
}
