package entities

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Author    string    `json:"author" db:"author"`
	Year      int       `json:"year" db:"year"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateBookDTO struct {
	Title  string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
	Year   int    `json:"year" validate:"required,min=1000,max=2034"`
}

type UpdateBookDTO struct {
	Title  string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
	Year   int    `json:"year" validate:"required,min=1000,max=2034"`
}

// Domain validation methods
func (dto *CreateBookDTO) Validate() error {
	if strings.TrimSpace(dto.Title) == "" {
		return ErrInvalidTitle
	}
	if strings.TrimSpace(dto.Author) == "" {
		return ErrInvalidAuthor
	}
	if dto.Year < 1000 || dto.Year > 2034 {
		return ErrInvalidYear
	}
	return nil
}

func (dto *UpdateBookDTO) Validate() error {
	if strings.TrimSpace(dto.Title) == "" {
		return ErrInvalidTitle
	}
	if strings.TrimSpace(dto.Author) == "" {
		return ErrInvalidAuthor
	}
	if dto.Year < 1000 || dto.Year > 2034 {
		return ErrInvalidYear
	}
	return nil
}

func (dto *CreateBookDTO) ToBook() *Book {
	return &Book{
		ID:     uuid.New(),
		Title:  strings.TrimSpace(dto.Title),
		Author: strings.TrimSpace(dto.Author),
		Year:   dto.Year,
	}
}

// ValidateBookData validates book data before persistence
func (b *Book) ValidateBookData() error {
	if b.Title == "" {
		return ErrInvalidTitle
	}
	if b.Author == "" {
		return ErrInvalidAuthor
	}
	if b.Year < 1000 || b.Year > 2034 {
		return ErrInvalidYear
	}
	return nil
}