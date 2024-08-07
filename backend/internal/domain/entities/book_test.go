package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateBookDTO_Validate(t *testing.T) {
	tests := []struct {
		name    string
		dto     CreateBookDTO
		wantErr bool
		errType error
	}{
		{
			name: "valid book",
			dto: CreateBookDTO{
				Title:  "The Go Programming Language",
				Author: "Alan Donovan",
				Year:   2015,
			},
			wantErr: false,
		},
		{
			name: "empty title",
			dto: CreateBookDTO{
				Title:  "",
				Author: "Alan Donovan",
				Year:   2015,
			},
			wantErr: true,
			errType: ErrInvalidTitle,
		},
		{
			name: "empty author",
			dto: CreateBookDTO{
				Title:  "The Go Programming Language",
				Author: "",
				Year:   2015,
			},
			wantErr: true,
			errType: ErrInvalidAuthor,
		},
		{
			name: "invalid year - too low",
			dto: CreateBookDTO{
				Title:  "The Go Programming Language",
				Author: "Alan Donovan",
				Year:   999,
			},
			wantErr: true,
			errType: ErrInvalidYear,
		},
		{
			name: "invalid year - too high",
			dto: CreateBookDTO{
				Title:  "The Go Programming Language",
				Author: "Alan Donovan",
				Year:   2035,
			},
			wantErr: true,
			errType: ErrInvalidYear,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dto.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateBookDTO_ToBook(t *testing.T) {
	dto := CreateBookDTO{
		Title:  "Clean Code",
		Author: "Robert Martin",
		Year:   2008,
	}

	book := dto.ToBook()

	assert.NotEmpty(t, book.ID)
	assert.Equal(t, "Clean Code", book.Title)
	assert.Equal(t, "Robert Martin", book.Author)
	assert.Equal(t, 2008, book.Year)
}

func TestBook_ValidateBookData(t *testing.T) {
	tests := []struct {
		name    string
		book    Book
		wantErr bool
		errType error
	}{
		{
			name: "valid book",
			book: Book{
				Title:  "Effective Go",
				Author: "The Go Team",
				Year:   2021,
			},
			wantErr: false,
		},
		{
			name: "empty title",
			book: Book{
				Title:  "",
				Author: "The Go Team",
				Year:   2021,
			},
			wantErr: true,
			errType: ErrInvalidTitle,
		},
		{
			name: "empty author",
			book: Book{
				Title:  "Effective Go",
				Author: "",
				Year:   2021,
			},
			wantErr: true,
			errType: ErrInvalidAuthor,
		},
		{
			name: "invalid year",
			book: Book{
				Title:  "Effective Go",
				Author: "The Go Team",
				Year:   500,
			},
			wantErr: true,
			errType: ErrInvalidYear,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.book.ValidateBookData()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}