package entities

import "errors"

var (
	ErrBookNotFound   = errors.New("book not found")
	ErrInvalidTitle   = errors.New("title cannot be empty")
	ErrInvalidAuthor  = errors.New("author cannot be empty")
	ErrInvalidYear    = errors.New("year must be between 1000 and 2034")
	ErrDatabaseError  = errors.New("database operation failed")
	ErrInvalidUUID    = errors.New("invalid UUID format")
)