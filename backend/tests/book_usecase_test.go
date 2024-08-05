package tests

import (
	"context"
	"testing"

	"byfood-library/internal/domain/entities"
	"byfood-library/internal/usecases"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// Mock BookRepository
type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) Create(ctx context.Context, book *entities.Book) (*entities.Book, error) {
	args := m.Called(ctx, book)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Book), args.Error(1)
}

func (m *MockBookRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Book, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Book), args.Error(1)
}

func (m *MockBookRepository) GetAll(ctx context.Context) ([]*entities.Book, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Book), args.Error(1)
}

func (m *MockBookRepository) Update(ctx context.Context, id uuid.UUID, book *entities.Book) (*entities.Book, error) {
	args := m.Called(ctx, id, book)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Book), args.Error(1)
}

func (m *MockBookRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupBookUseCase() (*MockBookRepository, usecases.BookUseCase) {
	mockRepo := new(MockBookRepository)
	logger, _ := zap.NewDevelopment()
	useCase := usecases.NewBookUseCase(mockRepo, logger)
	return mockRepo, useCase
}

func TestBookUseCase_CreateBook(t *testing.T) {
	mockRepo, useCase := setupBookUseCase()

	t.Run("successful creation", func(t *testing.T) {
		dto := &entities.CreateBookDTO{
			Title:  "The Go Programming Language",
			Author: "Alan Donovan",
			Year:   2015,
		}

		expectedBook := &entities.Book{
			ID:     uuid.New(),
			Title:  "The Go Programming Language",
			Author: "Alan Donovan",
			Year:   2015,
		}

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.Book")).Return(expectedBook, nil).Once()

		result, err := useCase.CreateBook(context.Background(), dto)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedBook.Title, result.Title)
		assert.Equal(t, expectedBook.Author, result.Author)
		assert.Equal(t, expectedBook.Year, result.Year)

		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		dto := &entities.CreateBookDTO{
			Title:  "", // Invalid title
			Author: "Alan Donovan",
			Year:   2015,
		}

		result, err := useCase.CreateBook(context.Background(), dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, entities.ErrInvalidTitle, err)

		// Repository should not be called for validation errors
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		dto := &entities.CreateBookDTO{
			Title:  "The Go Programming Language",
			Author: "Alan Donovan",
			Year:   2015,
		}

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.Book")).Return(nil, entities.ErrDatabaseError).Once()

		result, err := useCase.CreateBook(context.Background(), dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, entities.ErrDatabaseError, err)

		mockRepo.AssertExpectations(t)
	})
}

func TestBookUseCase_GetAllBooks(t *testing.T) {
	mockRepo, useCase := setupBookUseCase()

	t.Run("successful retrieval", func(t *testing.T) {
		expectedBooks := []*entities.Book{
			{ID: uuid.New(), Title: "Book 1", Author: "Author 1", Year: 2020},
			{ID: uuid.New(), Title: "Book 2", Author: "Author 2", Year: 2021},
		}

		mockRepo.On("GetAll", mock.Anything).Return(expectedBooks, nil).Once()

		result, err := useCase.GetAllBooks(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, expectedBooks[0].Title, result[0].Title)
		assert.Equal(t, expectedBooks[1].Title, result[1].Title)

		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo.On("GetAll", mock.Anything).Return(nil, entities.ErrDatabaseError).Once()

		result, err := useCase.GetAllBooks(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, entities.ErrDatabaseError, err)

		mockRepo.AssertExpectations(t)
	})
}

func TestBookUseCase_GetBookByID(t *testing.T) {
	mockRepo, useCase := setupBookUseCase()

	t.Run("successful retrieval", func(t *testing.T) {
		bookID := uuid.New()
		expectedBook := &entities.Book{
			ID:     bookID,
			Title:  "The Go Programming Language",
			Author: "Alan Donovan",
			Year:   2015,
		}

		mockRepo.On("GetByID", mock.Anything, bookID).Return(expectedBook, nil).Once()

		result, err := useCase.GetBookByID(context.Background(), bookID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedBook.ID, result.ID)
		assert.Equal(t, expectedBook.Title, result.Title)

		mockRepo.AssertExpectations(t)
	})

	t.Run("book not found", func(t *testing.T) {
		bookID := uuid.New()

		mockRepo.On("GetByID", mock.Anything, bookID).Return(nil, entities.ErrBookNotFound).Once()

		result, err := useCase.GetBookByID(context.Background(), bookID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, entities.ErrBookNotFound, err)

		mockRepo.AssertExpectations(t)
	})
}

func TestBookUseCase_UpdateBook(t *testing.T) {
	mockRepo, useCase := setupBookUseCase()

	t.Run("successful update", func(t *testing.T) {
		bookID := uuid.New()
		dto := &entities.UpdateBookDTO{
			Title:  "Updated Title",
			Author: "Updated Author",
			Year:   2022,
		}

		expectedBook := &entities.Book{
			ID:     bookID,
			Title:  "Updated Title",
			Author: "Updated Author",
			Year:   2022,
		}

		mockRepo.On("Update", mock.Anything, bookID, mock.AnythingOfType("*entities.Book")).Return(expectedBook, nil).Once()

		result, err := useCase.UpdateBook(context.Background(), bookID, dto)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedBook.Title, result.Title)
		assert.Equal(t, expectedBook.Author, result.Author)
		assert.Equal(t, expectedBook.Year, result.Year)

		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		bookID := uuid.New()
		dto := &entities.UpdateBookDTO{
			Title:  "", // Invalid title
			Author: "Updated Author",
			Year:   2022,
		}

		result, err := useCase.UpdateBook(context.Background(), bookID, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, entities.ErrInvalidTitle, err)

		// Repository should not be called for validation errors
		mockRepo.AssertExpectations(t)
	})
}

func TestBookUseCase_DeleteBook(t *testing.T) {
	mockRepo, useCase := setupBookUseCase()

	t.Run("successful deletion", func(t *testing.T) {
		bookID := uuid.New()

		mockRepo.On("Delete", mock.Anything, bookID).Return(nil).Once()

		err := useCase.DeleteBook(context.Background(), bookID)

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("book not found", func(t *testing.T) {
		bookID := uuid.New()

		mockRepo.On("Delete", mock.Anything, bookID).Return(entities.ErrBookNotFound).Once()

		err := useCase.DeleteBook(context.Background(), bookID)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrBookNotFound, err)

		mockRepo.AssertExpectations(t)
	})
}