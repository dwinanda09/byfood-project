package usecases

import (
	"context"
	"testing"
	"time"

	"byfood-library/internal/domain/entities"
	"byfood-library/internal/domain/repositories"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockBookRepository for testing
type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) Create(ctx context.Context, book *entities.Book) (*entities.Book, error) {
	args := m.Called(ctx, book)
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
	return args.Get(0).([]*entities.Book), args.Error(1)
}

func (m *MockBookRepository) Update(ctx context.Context, id uuid.UUID, book *entities.Book) (*entities.Book, error) {
	args := m.Called(ctx, id, book)
	return args.Get(0).(*entities.Book), args.Error(1)
}

func (m *MockBookRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupTest() (BookUseCase, *MockBookRepository) {
	mockRepo := new(MockBookRepository)
	logger := zap.NewNop()
	useCase := NewBookUseCase(mockRepo, logger)
	return useCase, mockRepo
}

func TestBookUseCase_CreateBook(t *testing.T) {
	useCase, mockRepo := setupTest()

	t.Run("successful creation", func(t *testing.T) {
		dto := &entities.CreateBookDTO{
			Title:  "Test Book",
			Author: "Test Author",
			Year:   2020,
		}

		expectedBook := &entities.Book{
			ID:        uuid.New(),
			Title:     dto.Title,
			Author:    dto.Author,
			Year:      dto.Year,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entities.Book")).Return(expectedBook, nil)

		result, err := useCase.CreateBook(context.Background(), dto)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, dto.Title, result.Title)
		assert.Equal(t, dto.Author, result.Author)
		assert.Equal(t, dto.Year, result.Year)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		dto := &entities.CreateBookDTO{
			Title:  "", // Invalid empty title
			Author: "Test Author",
			Year:   2020,
		}

		result, err := useCase.CreateBook(context.Background(), dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, entities.ErrInvalidTitle, err)
	})
}

func TestBookUseCase_GetBookByID(t *testing.T) {
	useCase, mockRepo := setupTest()

	t.Run("successful retrieval", func(t *testing.T) {
		bookID := uuid.New()
		expectedBook := &entities.Book{
			ID:     bookID,
			Title:  "Test Book",
			Author: "Test Author",
			Year:   2020,
		}

		mockRepo.On("GetByID", mock.Anything, bookID).Return(expectedBook, nil)

		result, err := useCase.GetBookByID(context.Background(), bookID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedBook.ID, result.ID)
		assert.Equal(t, expectedBook.Title, result.Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("book not found", func(t *testing.T) {
		bookID := uuid.New()

		mockRepo.On("GetByID", mock.Anything, bookID).Return(nil, entities.ErrBookNotFound)

		result, err := useCase.GetBookByID(context.Background(), bookID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, entities.ErrBookNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestBookUseCase_GetAllBooks(t *testing.T) {
	useCase, mockRepo := setupTest()

	t.Run("successful retrieval", func(t *testing.T) {
		expectedBooks := []*entities.Book{
			{
				ID:     uuid.New(),
				Title:  "Book 1",
				Author: "Author 1",
				Year:   2020,
			},
			{
				ID:     uuid.New(),
				Title:  "Book 2",
				Author: "Author 2",
				Year:   2021,
			},
		}

		mockRepo.On("GetAll", mock.Anything).Return(expectedBooks, nil)

		result, err := useCase.GetAllBooks(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, expectedBooks[0].Title, result[0].Title)
		assert.Equal(t, expectedBooks[1].Title, result[1].Title)
		mockRepo.AssertExpectations(t)
	})
}

func TestBookUseCase_UpdateBook(t *testing.T) {
	useCase, mockRepo := setupTest()

	t.Run("successful update", func(t *testing.T) {
		bookID := uuid.New()
		dto := &entities.UpdateBookDTO{
			Title:  "Updated Book",
			Author: "Updated Author",
			Year:   2022,
		}

		updatedBook := &entities.Book{
			ID:        bookID,
			Title:     dto.Title,
			Author:    dto.Author,
			Year:      dto.Year,
			UpdatedAt: time.Now(),
		}

		mockRepo.On("Update", mock.Anything, bookID, mock.AnythingOfType("*entities.Book")).Return(updatedBook, nil)

		result, err := useCase.UpdateBook(context.Background(), bookID, dto)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, dto.Title, result.Title)
		assert.Equal(t, dto.Author, result.Author)
		assert.Equal(t, dto.Year, result.Year)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		bookID := uuid.New()
		dto := &entities.UpdateBookDTO{
			Title:  "", // Invalid empty title
			Author: "Updated Author",
			Year:   2022,
		}

		result, err := useCase.UpdateBook(context.Background(), bookID, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, entities.ErrInvalidTitle, err)
	})
}

func TestBookUseCase_DeleteBook(t *testing.T) {
	useCase, mockRepo := setupTest()

	t.Run("successful deletion", func(t *testing.T) {
		bookID := uuid.New()

		mockRepo.On("Delete", mock.Anything, bookID).Return(nil)

		err := useCase.DeleteBook(context.Background(), bookID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("book not found", func(t *testing.T) {
		bookID := uuid.New()

		mockRepo.On("Delete", mock.Anything, bookID).Return(entities.ErrBookNotFound)

		err := useCase.DeleteBook(context.Background(), bookID)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrBookNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}