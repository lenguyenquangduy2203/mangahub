package users

import (
	"context"
	"errors"
	"log"
	manga "mangahub/internal/mangas"
	platform_errors "mangahub/pkg/errors"
	"mangahub/pkg/models"
)

// Define contract for accessing user library and progress data
type UserLibraryQuery interface {
	AddMangaToLibrary(ctx context.Context, userID string, mangaID string, currentChapter int) error
	GetLibrary(ctx context.Context, userID string, status string, limit int, offset int) (any, int64, error)
	UpdateReadingProgress(ctx context.Context, userID string, mangaID string, currentChapter int) error
}

// Define the business workflow contract
type UserLibraryService interface {
	AddMangaToUserLibrary(ctx context.Context, userID string, mangaID string, currentChapter int) error
	GetUserLibrary(ctx context.Context, userID string, status string, limit int, offset int) (any, error)
	UpdateUserReadingProgress(ctx context.Context, userID string, mangaID string, currentChapter int) error
}

type Service struct {
	UserLibraryQuery UserLibraryQuery
}

func NewService(uq UserLibraryQuery) *Service {
	return &Service{
		UserLibraryQuery: uq,
	}
}

func (s *Service) AddMangaToUserLibrary(ctx context.Context, userID string, mangaID string, currentChapter int) error {
	err := s.UserLibraryQuery.AddMangaToLibrary(ctx, userID, mangaID, currentChapter)
	if err != nil {
		if errors.Is(err, manga.ErrMangaNotExistInDatabase) || errors.Is(err, ErrMangaAlreadyInUserLibrary) || errors.Is(err, manga.ErrInValidCurrentChapter) {
			return err
		}
		return platform_errors.ErrDatabaseOperation
	}

	return nil
}

func (s *Service) UpdateUserReadingProgress(ctx context.Context, userID string, mangaID string, currentChapter int) error {
	err := s.UserLibraryQuery.UpdateReadingProgress(ctx, userID, mangaID, currentChapter)

	if err != nil {
		log.Println(err)

		if errors.Is(err, manga.ErrMangaNotExistInDatabase) || errors.Is(err, ErrMangaNotExistInUserLibrary) || errors.Is(err, manga.ErrInValidCurrentChapter) {
			return err
		}

		return platform_errors.ErrDatabaseOperation
	}

	return nil
}

func (s *Service) GetUserLibrary(ctx context.Context, userID string, status string, limit int, offset int) (any, error) {
	library, total, err := s.UserLibraryQuery.GetLibrary(ctx, userID, status, limit, offset)
	if err != nil {
		return nil, platform_errors.ErrDatabaseOperation
	}

	if lib, ok := library.([]models.UserLibrary); ok {
		paginatedLibrary := models.PaginatedUserLibrary{
			Total:   total,
			Limit:   limit,
			Offset:  offset,
			Results: _convertToLibraryItems(lib),
		}
		return paginatedLibrary, nil
	} else if lib, ok := library.(map[string][]models.UserLibrary); ok {
		readingList := models.ReadingList{
			Reading:    _convertToLibraryItems(lib["reading"]),
			Completed:  _convertToLibraryItems(lib["completed"]),
			PlanToRead: _convertToLibraryItems(lib["plan_to_read"]),
		}
		return readingList, nil
	}

	return nil, errors.New("unable to parse library data")
}

func _convertToLibraryItems(ul []models.UserLibrary) []models.LibraryItem {
	items := make([]models.LibraryItem, 0, len(ul))
	for _, item := range ul {
		items = append(items, models.LibraryItem{
			MangaID:        item.MangaID,
			CurrentChapter: item.CurrentChapter,
			Status:         item.Status,
			UpdatedAt:      item.UpdatedAt,
		})
	}
	return items
}
