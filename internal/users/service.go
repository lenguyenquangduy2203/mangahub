package users

import (
	"context"
	"errors"
	"log"
	manga "mangahub/internal/mangas"
	platform_errors "mangahub/pkg/errors"
)

// Define contract for accessing user library and progress data
type UserLibraryQuery interface {
	AddMangaToLibrary(ctx context.Context, userID string, mangaID string, currentChapter int) error
	GetLibrary(ctx context.Context, userID string, status string, limit int, offset int) (any, error)
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
		if errors.Is(err, manga.ErrMangaNotExistInDatabase) || errors.Is(err, ErrMangaAlreadyInUserLibrary) {
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
	return nil, nil
}
