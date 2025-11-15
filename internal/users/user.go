package users

import (
	"context"
	"errors"
	"mangahub/internal/mangas"
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

	if currentChapter > 0 {
		userManga.Status = enums.READING_READING.StringUpper()
	} else {
		userManga.Status = enums.READING_PLAN_TO_READ.StringUpper()
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

func (r *Repository) GetLibrary(ctx context.Context, userID string, status string, limit int, offset int) (any, error) {
	// 1. Start GORM query builder
	db := r.DB.DB().WithContext(ctx).Model(&models.UserLibrary{})

	// 2. Filter by UserID
	db = db.Where("user_id = ?", userID)

	// 3. Optional: Filter by Status if provided
	if status != "" {
		// You may need to normalize the status string here (e.g., to UPPERCASE)
		// based on how you store it in the database.
		db = db.Where("status = ?", strings.ToUpper(status))
	}

	// 4. Apply Pagination
	if limit > 0 {
		db = db.Limit(limit)
	}
	if offset > 0 {
		db = db.Offset(offset)
	}

	// 5. Execute the query
	var libraryEntries []models.UserLibrary
	// Use .Find() to retrieve multiple records
	if err := db.Find(&libraryEntries).Error; err != nil {
		return nil, err
	}

	// Since the interface returns 'any', return the list of library entries.
	return libraryEntries, nil
}
