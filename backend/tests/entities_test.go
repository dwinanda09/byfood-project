package tests

import (
	"testing"

	"byfood-library/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestCreateBookDTO_Validate(t *testing.T) {
	tests := []struct {
		name    string
		dto     entities.CreateBookDTO
		wantErr error
	}{
		{
			name: "valid book",
			dto: entities.CreateBookDTO{
				Title:  "The Go Programming Language",
				Author: "Alan Donovan",
				Year:   2015,
			},
			wantErr: nil,
		},
		{
			name: "empty title",
			dto: entities.CreateBookDTO{
				Title:  "",
				Author: "Alan Donovan",
				Year:   2015,
			},
			wantErr: entities.ErrInvalidTitle,
		},
		{
			name: "empty author",
			dto: entities.CreateBookDTO{
				Title:  "The Go Programming Language",
				Author: "",
				Year:   2015,
			},
			wantErr: entities.ErrInvalidAuthor,
		},
		{
			name: "invalid year too low",
			dto: entities.CreateBookDTO{
				Title:  "The Go Programming Language",
				Author: "Alan Donovan",
				Year:   999,
			},
			wantErr: entities.ErrInvalidYear,
		},
		{
			name: "invalid year too high",
			dto: entities.CreateBookDTO{
				Title:  "The Go Programming Language",
				Author: "Alan Donovan",
				Year:   2035,
			},
			wantErr: entities.ErrInvalidYear,
		},
		{
			name: "whitespace only title",
			dto: entities.CreateBookDTO{
				Title:  "   ",
				Author: "Alan Donovan",
				Year:   2015,
			},
			wantErr: entities.ErrInvalidTitle,
		},
		{
			name: "whitespace only author",
			dto: entities.CreateBookDTO{
				Title:  "The Go Programming Language",
				Author: "   ",
				Year:   2015,
			},
			wantErr: entities.ErrInvalidAuthor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dto.Validate()
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateBookDTO_Validate(t *testing.T) {
	tests := []struct {
		name    string
		dto     entities.UpdateBookDTO
		wantErr error
	}{
		{
			name: "valid book",
			dto: entities.UpdateBookDTO{
				Title:  "Clean Code",
				Author: "Robert Martin",
				Year:   2008,
			},
			wantErr: nil,
		},
		{
			name: "empty title",
			dto: entities.UpdateBookDTO{
				Title:  "",
				Author: "Robert Martin",
				Year:   2008,
			},
			wantErr: entities.ErrInvalidTitle,
		},
		{
			name: "empty author",
			dto: entities.UpdateBookDTO{
				Title:  "Clean Code",
				Author: "",
				Year:   2008,
			},
			wantErr: entities.ErrInvalidAuthor,
		},
		{
			name: "invalid year",
			dto: entities.UpdateBookDTO{
				Title:  "Clean Code",
				Author: "Robert Martin",
				Year:   999,
			},
			wantErr: entities.ErrInvalidYear,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dto.Validate()
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateBookDTO_ToBook(t *testing.T) {
	dto := entities.CreateBookDTO{
		Title:  "  The Go Programming Language  ",
		Author: "  Alan Donovan  ",
		Year:   2015,
	}

	book := dto.ToBook()

	assert.NotNil(t, book)
	assert.NotEqual(t, "", book.ID.String()) // UUID should be generated
	assert.Equal(t, "The Go Programming Language", book.Title) // Should be trimmed
	assert.Equal(t, "Alan Donovan", book.Author)              // Should be trimmed
	assert.Equal(t, 2015, book.Year)
}