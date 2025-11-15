package mangas

import (
	"errors"
	platform_errors "mangahub/pkg/errors"
)

var ErrInvalidMangaID = errors.New("invalid manga id")
var ErrInValidCurrentChapter = platform_errors.ErrInValidCurrentChapter // re-export from shared package
var ErrMangaNotExistInDatabase = errors.New("manga does not exist in database")
