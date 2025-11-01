package database

import (
	"log"
	"mangahub/data/models"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabaseConnection(dbPath string) (DBConnector, error) {

	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, err
	}

	gormDB, err := gorm.Open(sqlite.Open(absPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := autoMigrate(gormDB); err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Return the injectable struct
	return &Database{gormDB: gormDB, sqlDB: sqlDB}, nil
}

func autoMigrate(db *gorm.DB) error {
	log.Println("Database: Running AutoMigrate...")
	return db.AutoMigrate(
		&models.User{},
		&models.Manga{},
		&models.UserProgress{},
	)
}
