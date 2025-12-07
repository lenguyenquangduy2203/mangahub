package models

import (
	"errors"
	platform_errors "mangahub/pkg/errors"
	"mangahub/pkg/models/enums"
	gormHelper "mangahub/pkg/utils/gorm"
	"time"

	"gorm.io/gorm"
)

// Domain models
type User struct {
	ID           string    `gorm:"primaryKey;type:TEXT"`
	Username     string    `gorm:"unique;type:TEXT"`
	PasswordHash string    `gorm:"column:password_hash;type:TEXT"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}

type Manga struct {
	ID             string `gorm:"primaryKey;type:TEXT"`
	Title          string `gorm:"type:TEXT"`
	Author         string `gorm:"type:TEXT"`
	Genres         string `gorm:"type:TEXT"`
	Status         string `gorm:"type:TEXT"` // ONGOING, HIATUS, COMPLETED
	TotalChapters  int    `gorm:"column:total_chapters"`
	Description    string `gorm:"type:TEXT"`
	MangaUpdatesID *int   `gorm:"column:mangaupdates_id;uniqueIndex;type:INTEGER;default:null"`
}

type UserLibrary struct {
	// Composite Primary Key: (user_id, manga_id)
	UserID         string    `gorm:"primaryKey;type:TEXT"`
	MangaID        string    `gorm:"primaryKey;type:TEXT"`
	CurrentChapter int       `gorm:"column:current_chapter;default:0"`
	Status         string    `gorm:"type:TEXT"` // READING, COMPLETED, PLAN_TO_READ
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// Adding the 'foreignKey' tag for clarity.
	User  User  `gorm:"foreignKey:UserID;references:ID"`
	Manga Manga `gorm:"foreignKey:MangaID;references:ID"`
}

func (ul *UserLibrary) BeforeCreate(tx *gorm.DB) error {
	manga, err := gormHelper.FindOne[Manga](tx, tx.Statement.Context, "id = ?", ul.MangaID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return platform_errors.ErrMangaNotExistInDatabase
		}

		return err
	}

	if ul.CurrentChapter == manga.TotalChapters && manga.Status == enums.MANGA_COMPLETED.StringUpper() {
		ul.Status = enums.READING_COMPLETED.StringUpper()
	} else if ul.CurrentChapter < manga.TotalChapters && ul.CurrentChapter > 0 {
		ul.Status = enums.READING_READING.StringUpper()
	} else if ul.CurrentChapter == 0 {
		ul.Status = enums.READING_PLAN_TO_READ.StringUpper()
	} else {
		return platform_errors.ErrInValidCurrentChapter
	}

	return nil
}

func (ul *UserLibrary) BeforeUpdate(tx *gorm.DB) error {
	manga, err := gormHelper.FindOne[Manga](tx, tx.Statement.Context, "id = ?", ul.MangaID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return platform_errors.ErrMangaNotExistInDatabase
		}

		return err
	}

	if ul.CurrentChapter == manga.TotalChapters && manga.Status == enums.MANGA_COMPLETED.StringUpper() {
		ul.Status = enums.READING_COMPLETED.StringUpper()
	} else if ul.CurrentChapter < manga.TotalChapters && ul.CurrentChapter > 0 {
		ul.Status = enums.READING_READING.StringUpper()
	} else {
		return platform_errors.ErrInValidCurrentChapter
	}

	return nil
}

type MangaSearchQuery struct {
	Title  string
	Author string
	Genre  string
	Status string
	Limit  int
	Offset int
}

type MangaListItem struct {
	ID            string
	Title         string
	TotalChapters int
	Status        string
}

type PaginatedMangaResult struct {
	Total   int64
	Limit   int
	Offset  int
	Results []MangaListItem
}

type ReadingList struct {
	Reading    []LibraryItem
	Completed  []LibraryItem
	PlanToRead []LibraryItem
}

type LibraryItem struct {
	MangaID        string
	CurrentChapter int
	Status         string
	UpdatedAt      time.Time
}

type PaginatedUserLibrary struct {
	Total   int64
	Limit   int
	Offset  int
	Results []LibraryItem
}
