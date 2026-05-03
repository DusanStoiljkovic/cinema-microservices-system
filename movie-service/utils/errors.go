package utils

import "errors"

var (
	ErrInvalidInput       = errors.New("invalid input")
	ErrNotFound           = errors.New("resource not found")
	ErrConflict           = errors.New("resource conflict")
	ErrGenreNotFound      = errors.New("one or more genres do not exist")
	ErrRecordNotFound     = errors.New("movie does not exist")
	ErrRecordAlreadyExist = errors.New("record already exist")
)
