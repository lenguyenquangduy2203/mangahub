package handlers

import (
	"context"
	"errors"
	"log"
	"mangahub/internal/mangas"
	"mangahub/pkg/models/dtos"
	"mangahub/pkg/pagination"
	"mangahub/pkg/utils/parsers"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type MangaHandler struct {
	MangaService mangas.MangaService
}

func NewMangaHandler(ms mangas.MangaService) *MangaHandler {
	return &MangaHandler{MangaService: ms}
}

func (h *MangaHandler) GetMangaV1(ctx *gin.Context) {
	var id dtos.MangaIDQuery

	if err := ctx.ShouldBindUri(&id); err != nil {
		errResp := dtos.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Wrong manga id format",
			Details: parsers.ParseValidationErrors(err),
		}

		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	goCtx := ctx.Request.Context()
	c, cancel := context.WithTimeout(goCtx, 3*time.Second)
	defer cancel()

	resp, err := h.MangaService.Get(c, id.ID)

	if err != nil {
		log.Printf("ERROR: during request for manga %s: %v", id.ID, err)

		if errors.Is(err, mangas.ErrInvalidMangaID) {
			ctx.JSON(http.StatusNotFound, dtos.ErrorResponse{
				Code:    "NOT_FOUND",
				Message: "Not found manga with id: " + id.ID,
			})

			return
		}

		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Get manga details failed due to internal error.",
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *MangaHandler) FindMangasV1(ctx *gin.Context) {
	query := dtos.MangaSearchQuery{
		Limit: pagination.DEFAULT_LIMIT, // Default limit
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

	log.Printf("DEBUG QUERY: %+v\n", query)

	goCtx := ctx.Request.Context()
	c, cancel := context.WithTimeout(goCtx, 3*time.Second)
	defer cancel()

	resp, err := h.MangaService.Find(c, query)

	if err != nil {
		log.Printf("ERROR: during request for finding mangas with query %v: %v", query, err)

		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Find mangas by query failed due to internal error.",
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
