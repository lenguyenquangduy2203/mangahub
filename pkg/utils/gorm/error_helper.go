package gorm

import (
	"errors"
	"strings"

	"github.com/mattn/go-sqlite3"
)

func IsUniqueConstraintError(err error) bool {
	var sqliteErr sqlite3.Error

	// Use errors.As to unwrap the GORM error and check if it contains
	// the underlying sqlite3.Error struct
	if errors.As(err, &sqliteErr) {
		// Check for the general constraint error code (19)
		// and if the error contains UNIQUE
		if errMsg := strings.ToUpper(sqliteErr.Error()); sqliteErr.Code == sqlite3.ErrConstraint && strings.Contains(errMsg, "UNIQUE") {
			return true
		}
	}
	return false
}
