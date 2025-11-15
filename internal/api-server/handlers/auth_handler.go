package handlers

import (
	"context"
	"errors"
	"log"
	"mangahub/internal/auth"
	"mangahub/pkg/models/dtos"
	"mangahub/pkg/utils/parsers"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService auth.AuthService
}

func NewAuthHandler(as auth.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: as}
}

func (h *AuthHandler) RegisterV1(ctx *gin.Context) {
	var body dtos.UserAuthRequest

	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		errResp := dtos.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Validation failed",
			Details: parsers.ParseValidationErrors(err),
		}

		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	goCtx := ctx.Request.Context()
	c, cancel := context.WithTimeout(goCtx, 3*time.Second)
	defer cancel()

	tokenInfo, err := h.AuthService.Register(c, body.Username, body.Password)

	if err != nil {
		log.Printf("ERROR: during registration for user %s: %v", body.Username, err)

		if errors.Is(err, auth.ErrUserConflict) {
			ctx.JSON(http.StatusConflict, dtos.ErrorResponse{
				Code:    "RESOURCE_CONFLICT",
				Message: "The requested username is already in use",
				Details: []dtos.ErrorDetail{
					{
						Field: "username",
						Issue: "Username already exists in the database",
					},
				},
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Registration failed due to internal error.",
		})
		return
	}

	resp := dtos.TokenResponse{
		AccessToken: tokenInfo.AccessToken,
		TokenType:   tokenInfo.TokenType,
		ExpiresAt:   tokenInfo.ExpiresAt,
	}
	ctx.JSON(http.StatusOK, resp)
}

func (h *AuthHandler) LoginV1(ctx *gin.Context) {
	var body dtos.UserAuthRequest

	if err := ctx.ShouldBindBodyWithJSON(&body); err != nil {
		errResp := dtos.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Validation failed",
			Details: parsers.ParseValidationErrors(err),
		}

		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	goCtx := ctx.Request.Context()
	c, cancel := context.WithTimeout(goCtx, 3*time.Second)
	defer cancel()

	tokenInfo, err := h.AuthService.Login(c, body.Username, body.Password)

	if err != nil {
		log.Printf("ERROR: during login request for user %s: %v", body.Username, err)

		if errors.Is(err, auth.ErrInvalidCredentials) {
			ctx.JSON(http.StatusUnauthorized, dtos.ErrorResponse{
				Code:    "UNAUTHORIZED",
				Message: "Wrong credentials",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Login failed due to internal error.",
		})
		return
	}

	resp := dtos.TokenResponse{
		AccessToken: tokenInfo.AccessToken,
		TokenType:   tokenInfo.TokenType,
		ExpiresAt:   tokenInfo.ExpiresAt,
	}
	ctx.JSON(http.StatusOK, resp)
}
