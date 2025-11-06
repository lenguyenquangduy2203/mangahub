package database

import (
	"database/sql"

	"gorm.io/gorm"
)

type DBConnector interface {
	DB() *gorm.DB
	SQLDB() *sql.DB
}

type Database struct {
	gormDB *gorm.DB
	sqlDB  *sql.DB
}

func (d *Database) DB() *gorm.DB {
	return d.gormDB
}

func (d *Database) SQLDB() *sql.DB {
	return d.sqlDB
}
