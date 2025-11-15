package handlers

import (
	"context"
	"errors"
	"log"
	"mangahub/internal/auth"
	"mangahub/internal/mangas"
	"mangahub/internal/users"
	"mangahub/pkg/models/dtos"
	"mangahub/pkg/utils/parsers"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserLibraryService users.UserLibraryService
}

func NewUserHandler(ls users.UserLibraryService) *UserHandler {
	return &UserHandler{UserLibraryService: ls}
}

func (h *UserHandler) AddMangaV1(ctx *gin.Context) {
	var body dtos.UserMangaAddRequest
	userID, err := auth.GetUserIDFromContext(ctx)

	if err != nil {
		errResp := dtos.ErrorResponse{
			Code:    "AUTH_ERROR",
			Message: "Could not retrieve authenticated user info",
		}

		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}

	if err = ctx.ShouldBindBodyWithJSON(&body); err != nil {
		errResp := dtos.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Wrong format",
			Details: parsers.ParseValidationErrors(err),
		}

		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	goCtx := ctx.Request.Context()
	c, cancel := context.WithTimeout(goCtx, 3*time.Second)
	defer cancel()

	err = h.UserLibraryService.AddMangaToUserLibrary(c, userID, body.MangaID, body.CurrentChapter)

	if err != nil {
		log.Printf("ERROR: during request for adding manga to library for %v: %v", userID, err)

		if errors.Is(err, users.ErrMangaAlreadyInUserLibrary) {
			errResp := dtos.ErrorResponse{
				Code:    "MANGA_ALREADY_IN_LIBRARY",
				Message: "Manga is already in user library",
			}

			ctx.JSON(http.StatusConflict, errResp)
			return
		}

		if errors.Is(err, mangas.ErrMangaNotExistInDatabase) {
			errResp := dtos.ErrorResponse{
				Code:    "MANGA_NOT_FOUND",
				Message: "Manga does not exist in database",
			}

			ctx.JSON(http.StatusNotFound, errResp)
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Add manga to user library failed due to internal error.",
		})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (h *UserHandler) UpdateReadingProgressV1(ctx *gin.Context) {
	var body dtos.UserMangaUpdateRequest
	userID, err := auth.GetUserIDFromContext(ctx)

	if err != nil {
		errResp := dtos.ErrorResponse{
			Code:    "AUTH_ERROR",
			Message: "Could not retrieve authenticated user info",
		}
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}

	if err = ctx.ShouldBindBodyWithJSON(&body); err != nil {
		errResp := dtos.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Wrong format",
			Details: parsers.ParseValidationErrors(err),
		}
		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}
	goCtx := ctx.Request.Context()
	c, cancel := context.WithTimeout(goCtx, 3*time.Second)
	defer cancel()

	err = h.UserLibraryService.UpdateUserReadingProgress(c, userID, body.MangaID, body.CurrentChapter)

	if err != nil {
		log.Printf("ERROR: during request for updating reading progress for %v: %v", userID, err)

		if errors.Is(err, users.ErrMangaNotExistInUserLibrary) {
			errResp := dtos.ErrorResponse{
				Code:    "MANGA_NOT_IN_LIBRARY",
				Message: "Manga does not exist in user library",
			}
			ctx.JSON(http.StatusNotFound, errResp)
			return
		}

		if errors.Is(err, mangas.ErrMangaNotExistInDatabase) {
			errResp := dtos.ErrorResponse{
				Code:    "MANGA_NOT_FOUND",
				Message: "Manga does not exist in database",
			}
			ctx.JSON(http.StatusNotFound, errResp)
			return
		}

		if errors.Is(err, mangas.ErrInValidCurrentChapter) {
			errResp := dtos.ErrorResponse{
				Code:    "INVALID_CURRENT_CHAPTER",
				Message: "Current chapter is invalid",
			}
			ctx.JSON(http.StatusBadRequest, errResp)
			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Update reading progress failed due to internal error.",
		})
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
