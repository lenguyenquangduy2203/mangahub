package mangas

import (
	"context"
	"errors"
	platform_errors "mangahub/pkg/errors"
	"mangahub/pkg/models"

	"gorm.io/gorm"
)

// Defines the contract for accessing manga data query
type MangaDataQuery interface {
	FindMangaByQuery(ctx context.Context, query models.MangaSearchQuery) ([]models.Manga, int64, error)
	GetMangaById(ctx context.Context, mangaID string) (*models.Manga, error)
}

// Defines the business logic workflow contract
type MangaService interface {
	Find(ctx context.Context, query models.MangaSearchQuery) (models.PaginatedMangaResult, error)
	Get(ctx context.Context, mangaID string) (models.Manga, error)
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

func (s *Service) Find(ctx context.Context, query models.MangaSearchQuery) (models.PaginatedMangaResult, error) {
	mangas, total, err := s.MangaQuery.FindMangaByQuery(ctx, query)
	if err != nil {
		return models.PaginatedMangaResult{}, platform_errors.ErrDatabaseOperation
	}

	results := make([]models.MangaListItem, 0, len(mangas))
	for _, m := range mangas {
		results = append(results, models.MangaListItem{
			ID:            m.ID,
			Title:         m.Title,
			TotalChapters: m.TotalChapters,
			Status:        m.Status,
		})
	}

	return models.PaginatedMangaResult{
		Total:   total,
		Limit:   query.Limit,
		Offset:  query.Offset,
		Results: results,
	}, nil
}

func (s *Service) Get(ctx context.Context, mangaID string) (models.Manga, error) {
	manga, err := s.MangaQuery.GetMangaById(ctx, mangaID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Manga{}, ErrInvalidMangaID
		}

		return models.Manga{}, platform_errors.ErrDatabaseOperation
	}

	return models.Manga{
		ID:            manga.ID,
		Title:         manga.Title,
		Author:        manga.Author,
		Genres:        manga.Genres,
		Status:        manga.Status,
		TotalChapters: manga.TotalChapters,
		Description:   manga.Description,
	}, nil
}
