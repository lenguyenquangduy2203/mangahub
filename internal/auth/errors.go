package auth

import "errors"

var ErrInvalidCredentials = errors.New("invalid username or password")

var ErrUserConflict = errors.New("username already exists")

var ErrTokenGeneration = errors.New("failed to generate security token")
