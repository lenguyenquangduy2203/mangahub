package models

import "time"

type User struct {
	ID           string    `gorm:"primaryKey;type:TEXT"`
	Username     string    `gorm:"unique;type:TEXT"`
	PasswordHash string    `gorm:"column:password_hash;type:TEXT"`
	CreatedAt    time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

type Manga struct {
	ID            string `gorm:"primaryKey;type:TEXT"`
	Title         string `gorm:"type:TEXT"`
	Author        string `gorm:"type:TEXT"`
	Genres        string `gorm:"type:TEXT"`
	Status        string `gorm:"type:TEXT"` // ONGOING, HIATUS, COMPLETE
	TotalChapters int    `gorm:"column:total_chapters"`
	Description   string `gorm:"type:TEXT"`
}

type UserProgress struct {
	// Composite Primary Key: (user_id, manga_id)
	UserID         string    `gorm:"primaryKey;type:TEXT"`
	MangaID        string    `gorm:"primaryKey;type:TEXT"`
	CurrentChapter int       `gorm:"column:current_chapter"`
	Status         string    `gorm:"type:TEXT"` // READING, COMPLETED, PLAN_TO_READ
	UpdatedAt      time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`

	// Adding the 'foreignKey' tag for clarity.
	User  User  `gorm:"foreignKey:UserID;references:ID"`
	Manga Manga `gorm:"foreignKey:MangaID;references:ID"`
}
