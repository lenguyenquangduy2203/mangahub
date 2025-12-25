package handlers

import (
	"context"
	"errors"
	"log"
	mangas2 "mangahub/internal/api-server/mangas"
	"mangahub/pkg/models"
	"mangahub/pkg/models/dtos"
	"mangahub/pkg/pagination"
	"mangahub/pkg/utils/parsers"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type MangaHandler struct {
	MangaService mangas2.MangaService
}

func NewMangaHandler(ms mangas2.MangaService) *MangaHandler {
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

	result, err := h.MangaService.Get(c, id.ID)
	resp := _mapMangaToDTO(result)

	if err != nil {
		log.Printf("ERROR: during request for manga %s: %v", id.ID, err)

		if errors.Is(err, mangas2.ErrInvalidMangaID) {
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
	queryDto := dtos.MangaSearchQuery{}

	if err := ctx.ShouldBindQuery(&queryDto); err != nil {
		errResp := dtos.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Wrong query format",
			Details: parsers.ParseValidationErrors(err),
		}

		ctx.JSON(http.StatusBadRequest, errResp)
		return
	}

	// Double check limit and page
	if queryDto.Limit <= 0 {
		queryDto.Limit = pagination.DEFAULT_LIMIT
	}
	if queryDto.Limit > pagination.MAX_LIMIT {
		queryDto.Limit = pagination.MAX_LIMIT
	}
	if queryDto.Page <= 0 {
		queryDto.Page = 1
	}

	offset := pagination.CalculateOffset(queryDto.Page, queryDto.Limit)

	log.Printf("DEBUG QUERY: %+v\n", queryDto)

	goCtx := ctx.Request.Context()
	c, cancel := context.WithTimeout(goCtx, 3*time.Second)
	defer cancel()

	query := models.MangaSearchQuery{
		Title:  queryDto.Title,
		Author: queryDto.Author,
		Genre:  queryDto.Genre,
		Status: strings.ToUpper(queryDto.Status),
		Limit:  queryDto.Limit,
		Offset: offset,
	}

	result, err := h.MangaService.Find(c, query)

	if err != nil {
		log.Printf("ERROR: during request for finding mangas with query %v: %v", query, err)

		ctx.JSON(http.StatusInternalServerError, dtos.ErrorResponse{
			Code:    "INTERNAL_ERROR",
			Message: "Find mangas by query failed due to internal error.",
		})
		return
	}

	// Use generic function to build public paginated response
	resp := dtos.BuildPaginatedResponse[dtos.MangaListItem](_mapPaginatedMangaResultToDTO(result))

	ctx.JSON(http.StatusOK, resp)
}

func _splitGenres(s string) []string {
	raw := strings.Split(s, ",")
	for i := range raw {
		raw[i] = strings.TrimSpace(raw[i])
	}
	return raw
}

func _mapMangaToDTO(manga models.Manga) dtos.MangaDetail {
	return dtos.MangaDetail{
		MangaID:       manga.ID,
		Title:         manga.Title,
		Author:        manga.Author,
		Genres:        _splitGenres(manga.Genres),
		Status:        manga.Status,
		TotalChapters: manga.TotalChapters,
		Description:   manga.Description,
	}
}

func _mapMangaListItemToDTO(item models.MangaListItem) dtos.MangaListItem {
	return dtos.MangaListItem{
		MangaID:       item.ID,
		Title:         item.Title,
		TotalChapters: item.TotalChapters,
		Status:        item.Status,
	}
}

func _mapPaginatedMangaResultToDTO(result models.PaginatedMangaResult) dtos.PaginatedMangaList {
	items := make([]dtos.MangaListItem, 0, len(result.Results))
	for _, item := range result.Results {
		items = append(items, _mapMangaListItemToDTO(item))
	}

	return dtos.PaginatedMangaList{
		Total:   int(result.Total),
		Limit:   result.Limit,
		Offset:  result.Offset,
		Results: items,
	}
}
