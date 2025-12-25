package platform_errors

import "errors"

var ErrDatabaseOperation = errors.New("database operation failed")
var ErrInValidCurrentChapter = errors.New("invalid current chapter")
var ErrMangaNotExistInDatabase = errors.New("manga does not exist in database")
