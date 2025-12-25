package database

import (
	"fmt"
	"log"
	"mangahub/pkg/models"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewDatabaseConnection(dbPath string) (DBConnector, error) {

	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, err
	}

	gormDB, err := gorm.Open(sqlite.Open(absPath), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}

	// Enable FK
	if err := gormDB.Exec("PRAGMA foreign_keys = ON;").Error; err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Verify FK are enabled
	var result int
	gormDB.Raw("SELECT foreign_keys FROM pragma_foreign_keys();").Scan(&result)
	log.Printf("Foreign keys status: %d (1=enabled, 0=disabled)", result)

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
		&models.UserLibrary{},
	)
}
