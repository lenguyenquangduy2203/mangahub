package manga

import (
	"context"
	"errors"
	platform_errors "mangahub/pkg/errors"
	"mangahub/pkg/models"
	"mangahub/pkg/models/dtos"
	"strings"

	"gorm.io/gorm"
)

// Defines the contract for accessing manga data query
type MangaDataQuery interface {
	FindMangaByQuery(ctx context.Context, query dtos.MangaSearchQuery) ([]models.Manga, int64, error)
	GetMangaById(ctx context.Context, mangaID string) (*models.Manga, error)
}

// Defines the business logic workflow contract
type MangaService interface {
	Find(ctx context.Context, query dtos.MangaSearchQuery) (dtos.PaginatedMangaList, error)
	Get(ctx context.Context, mangaID string) (dtos.MangaDetail, error)
}

type Service struct {
	MangaQuery MangaDataQuery
}

// Constructor for the MangaService
func NewService(mq MangaDataQuery) *Service {
	return &Service{
		MangaQuery: mq,
	}
}

func (s *Service) Find(ctx context.Context, query dtos.MangaSearchQuery) (dtos.PaginatedMangaList, error) {
	mangas, total, err := s.MangaQuery.FindMangaByQuery(ctx, query)
	if err != nil {
		return dtos.PaginatedMangaList{}, platform_errors.ErrDatabaseOperation
	}

	results := make([]dtos.MangaListItem, 0, len(mangas))
	for _, m := range mangas {
		results = append(results, dtos.MangaListItem{
			MangaID:       m.ID,
			Title:         m.Title,
			TotalChapters: m.TotalChapters,
			Status:        m.Status,
		})
	}

	return dtos.PaginatedMangaList{
		Total:   int(total),
		Limit:   query.Limit,
		Offset:  query.Offset,
		Results: results,
	}, nil
}

func (s *Service) Get(ctx context.Context, mangaID string) (dtos.MangaDetail, error) {
	manga, err := s.MangaQuery.GetMangaById(ctx, mangaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dtos.MangaDetail{}, ErrInvalidMangaID
		}

		return dtos.MangaDetail{}, platform_errors.ErrDatabaseOperation
	}

	return dtos.MangaDetail{
		MangaID:       manga.ID,
		Title:         manga.Title,
		Author:        manga.Author,
		Genres:        splitGenres(manga.Genres),
		Status:        manga.Status,
		TotalChapters: manga.TotalChapters,
		Description:   manga.Description,
	}, nil
}

func splitGenres(s string) []string {
	raw := strings.Split(s, ",")
	for i := range raw {
		raw[i] = strings.TrimSpace(raw[i])
	}
	return raw
}
