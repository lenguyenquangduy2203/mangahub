package mangas

import "errors"

var ErrInvalidMangaID = errors.New("invalid manga id")
var ErrInValidCurrentChapter = errors.New("invalid current chapter")
var ErrMangaNotExistInDatabase = errors.New("manga does not exist in database")
