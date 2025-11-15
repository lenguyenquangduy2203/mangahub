package handlers

import (
	"context"
	"errors"
	"log"
	"mangahub/internal/auth"
	"mangahub/internal/mangas"
	"mangahub/internal/users"
	"mangahub/pkg/models"
	"mangahub/pkg/models/dtos"
	"mangahub/pkg/pagination"
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

func (h *UserHandler) GetUserLibraryV1(ctx *gin.Context) {
	userID, err := auth.GetUserIDFromContext(ctx)

	if err != nil {
		errResp := dtos.ErrorResponse{
			Code:    "AUTH_ERROR",
			Message: "Could not retrieve authenticated user info",
		}
		ctx.JSON(http.StatusInternalServerError, errResp)
		return
	}

	query := dtos.UserMangaGetRequest{
		Limit: pagination.DEFAULT_LIMIT,
	}

	if err := ctx.ShouldBindQuery(&query); err != nil {
		errResp := dtos.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Wrong query format",
			Details: parsers.ParseValidationErrors(err),
		}

		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	if query.Limit > pagination.MAX_LIMIT {
		query.Limit = pagination.MAX_LIMIT
	}

	goCtx := ctx.Request.Context()
	c, cancel := context.WithTimeout(goCtx, 3*time.Second)
	defer cancel()

	library, err := h.UserLibraryService.GetUserLibrary(c, userID, query.Status, query.Limit, query.Offset)

	resp := _mapLibraryToDTO(library)

	if err != nil || resp == nil {
		log.Printf("ERROR: during request for getting user library for %v: %v", userID, err)

		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Get user library failed due to internal error.",
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func _mapLibraryItemsToDTO(items []models.LibraryItem) []dtos.UserLibraryItem {
	result := make([]dtos.UserLibraryItem, 0, len(items))
	for _, item := range items {
		result = append(result, dtos.UserLibraryItem{
			MangaID:        item.MangaID,
			CurrentChapter: item.CurrentChapter,
			Status:         item.Status,
			LastUpdated:    item.UpdatedAt,
		})
	}
	return result
}

func _mapLibraryToDTO(library any) any {
	if lib, ok := library.(models.PaginatedUserLibrary); ok {
		resp := dtos.PaginatedUserLibrary{
			Total:   int(lib.Total),
			Limit:   lib.Limit,
			Offset:  lib.Offset,
			Results: _mapLibraryItemsToDTO(lib.Results),
		}

		return resp
	} else if lib, ok := library.(models.ReadingList); ok {
		resp := dtos.UserLibrary{
			ReadingList: dtos.ReadingList{
				Reading:    _mapLibraryItemsToDTO(lib.Reading),
				Completed:  _mapLibraryItemsToDTO(lib.Completed),
				PlanToRead: _mapLibraryItemsToDTO(lib.PlanToRead),
			},
		}

		return resp
	}

	return nil
}
