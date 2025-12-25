package users

import "errors"

var ErrMangaAlreadyInUserLibrary = errors.New("manga already exists in user library")
var ErrMangaNotExistInUserLibrary = errors.New("manga not exists in user library")
