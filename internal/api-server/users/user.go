package users

import (
	"context"
	"errors"
	"log"
	"mangahub/internal/api-server/mangas"
	"mangahub/pkg/models"
	"mangahub/pkg/models/enums"
	"mangahub/pkg/utils/database"
	gormHelper "mangahub/pkg/utils/gorm"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	DB database.DBConnector
}

func NewRepository(db database.DBConnector) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) CreateUser(ctx context.Context, username, hashedPassword string) (*models.User, error) {
	id := uuid.New().String()

	user := models.User{
		ID:           id,
		Username:     username,
		PasswordHash: hashedPassword,
	}

	err := gormHelper.Create(r.DB.DB(), ctx, &user)

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.New("duplicate uuid")
		}

		return nil, err
	}

	return &user, nil
}

func (r *Repository) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := gormHelper.FindOne[models.User](r.DB.DB(), ctx, "username = ?", username)

	if err != nil {
		return nil, err
	}

	return &user, err
}

func (r *Repository) AddMangaToLibrary(ctx context.Context, userID string, mangaID string, currentChapter int) error {
	userManga := models.UserLibrary{
		UserID:         userID,
		MangaID:        mangaID,
		CurrentChapter: currentChapter,
	}

	err := gormHelper.Create(r.DB.DB(), ctx, &userManga)

	if err != nil {
		if gormHelper.IsForeignKeyConstraintError(err) {
			return mangas.ErrMangaNotExistInDatabase
		}
		if gormHelper.IsUniqueConstraintError(err) {
			return ErrMangaAlreadyInUserLibrary
		}
	}
	log.Println(err)
	return err
}

func (r *Repository) UpdateReadingProgress(ctx context.Context, userID string, mangaID string, currentChapter int) error {
	userManga := models.UserLibrary{
		UserID:         userID,
		MangaID:        mangaID,
		CurrentChapter: currentChapter,
	}

	rows, err := gormHelper.Update(r.DB.DB(), ctx,
		&userManga,
		"user_id = ? AND manga_id = ?", userID, mangaID)

	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrMangaNotExistInUserLibrary
	}

	return nil
}

func (r *Repository) GetLibrary(ctx context.Context, userID string, status string, limit int, offset int) (any, int64, error) {
	var total int64

	tx := r.DB.DB().WithContext(ctx).Model(&models.UserLibrary{})

	tx = tx.Where("user_id = ?", userID)

	if status != "" {
		tx = tx.Where("status = ?", strings.ToUpper(status))

		// Count total before pagination
		if err := tx.Count(&total).Error; err != nil {
			return nil, 0, err
		}

		if limit > 0 {
			tx = tx.Limit(limit)
		}
		if offset > 0 {
			tx = tx.Offset(offset)
		}

		var libraryEntries []models.UserLibrary
		if err := tx.Preload("Manga").Order("manga_id ASC").Find(&libraryEntries).Error; err != nil {
			return nil, 0, err
		}

		return libraryEntries, total, nil
	} else {
		var readingLibraryEntries []models.UserLibrary
		var completedLibraryEntries []models.UserLibrary
		var planToReadLibraryEntries []models.UserLibrary

		if err := r.DB.DB().WithContext(ctx).Model(&models.UserLibrary{}).
			Preload("Manga").
			Where("user_id = ? AND status = ?", userID, enums.READING_READING.StringUpper()).
			Order("manga_id ASC").
			Limit(5).
			Find(&readingLibraryEntries).Error; err != nil {
			return nil, 0, err
		}

		if err := r.DB.DB().WithContext(ctx).Model(&models.UserLibrary{}).
			Preload("Manga").
			Where("user_id = ? AND status = ?", userID, enums.READING_COMPLETED.StringUpper()).
			Order("manga_id ASC").
			Limit(5).
			Find(&completedLibraryEntries).Error; err != nil {
			return nil, 0, err
		}

		if err := r.DB.DB().WithContext(ctx).Model(&models.UserLibrary{}).
			Preload("Manga").
			Where("user_id = ? AND status = ?", userID, enums.READING_PLAN_TO_READ.StringUpper()).
			Order("manga_id ASC").
			Limit(5).
			Find(&planToReadLibraryEntries).Error; err != nil {
			return nil, 0, err
		}

		readingList := map[string][]models.UserLibrary{
			"reading":      readingLibraryEntries,
			"completed":    completedLibraryEntries,
			"plan_to_read": planToReadLibraryEntries,
		}

		return readingList, 0, nil
	}
}
