package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"byfood-library/internal/delivery/http/handlers"
	"byfood-library/internal/domain/entities"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// Mock BookUseCase
type MockBookUseCase struct {
	mock.Mock
}

func (m *MockBookUseCase) CreateBook(ctx context.Context, dto *entities.CreateBookDTO) (*entities.Book, error) {
	args := m.Called(ctx, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Book), args.Error(1)
}

func (m *MockBookUseCase) GetBookByID(ctx context.Context, id uuid.UUID) (*entities.Book, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Book), args.Error(1)
}

func (m *MockBookUseCase) GetAllBooks(ctx context.Context) ([]*entities.Book, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Book), args.Error(1)
}

func (m *MockBookUseCase) UpdateBook(ctx context.Context, id uuid.UUID, dto *entities.UpdateBookDTO) (*entities.Book, error) {
	args := m.Called(ctx, id, dto)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Book), args.Error(1)
}

func (m *MockBookUseCase) DeleteBook(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupBookHandler() (*MockBookUseCase, handlers.BookHandlerInterface) {
	mockUseCase := new(MockBookUseCase)
	logger, _ := zap.NewDevelopment()
	handler := handlers.NewBookHandler(mockUseCase, logger)
	return mockUseCase, handler
}

func TestBookHandler_GetBooks(t *testing.T) {
	mockUseCase, handler := setupBookHandler()

	t.Run("successful retrieval", func(t *testing.T) {
		expectedBooks := []*entities.Book{
			{ID: uuid.New(), Title: "Book 1", Author: "Author 1", Year: 2020},
			{ID: uuid.New(), Title: "Book 2", Author: "Author 2", Year: 2021},
		}

		mockUseCase.On("GetAllBooks", mock.Anything).Return(expectedBooks, nil).Once()

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/books", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetBooks(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var result []*entities.Book
		err = json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, expectedBooks[0].Title, result[0].Title)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("use case error", func(t *testing.T) {
		mockUseCase.On("GetAllBooks", mock.Anything).Return(nil, entities.ErrDatabaseError).Once()

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/books", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.GetBooks(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var errorResp handlers.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &errorResp)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to retrieve books", errorResp.Error)

		mockUseCase.AssertExpectations(t)
	})
}

func TestBookHandler_GetBook(t *testing.T) {
	mockUseCase, handler := setupBookHandler()

	t.Run("successful retrieval", func(t *testing.T) {
		bookID := uuid.New()
		expectedBook := &entities.Book{
			ID:     bookID,
			Title:  "The Go Programming Language",
			Author: "Alan Donovan",
			Year:   2015,
		}

		mockUseCase.On("GetBookByID", mock.Anything, bookID).Return(expectedBook, nil).Once()

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/books/"+bookID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(bookID.String())

		err := handler.GetBook(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var result entities.Book
		err = json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, expectedBook.ID, result.ID)
		assert.Equal(t, expectedBook.Title, result.Title)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("invalid UUID", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/books/invalid-uuid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid-uuid")

		err := handler.GetBook(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var errorResp handlers.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &errorResp)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid UUID format", errorResp.Error)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("book not found", func(t *testing.T) {
		bookID := uuid.New()

		mockUseCase.On("GetBookByID", mock.Anything, bookID).Return(nil, entities.ErrBookNotFound).Once()

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/books/"+bookID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(bookID.String())

		err := handler.GetBook(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var errorResp handlers.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &errorResp)
		assert.NoError(t, err)
		assert.Equal(t, "Book not found", errorResp.Error)

		mockUseCase.AssertExpectations(t)
	})
}

func TestBookHandler_CreateBook(t *testing.T) {
	mockUseCase, handler := setupBookHandler()

	t.Run("successful creation", func(t *testing.T) {
		dto := entities.CreateBookDTO{
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

		mockUseCase.On("CreateBook", mock.Anything, &dto).Return(expectedBook, nil).Once()

		reqBody, _ := json.Marshal(dto)
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateBook(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var result entities.Book
		err = json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, expectedBook.Title, result.Title)
		assert.Equal(t, expectedBook.Author, result.Author)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateBook(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var errorResp handlers.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &errorResp)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid request body", errorResp.Error)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		dto := entities.CreateBookDTO{
			Title:  "", // Invalid
			Author: "Alan Donovan",
			Year:   2015,
		}

		mockUseCase.On("CreateBook", mock.Anything, &dto).Return(nil, entities.ErrInvalidTitle).Once()

		reqBody, _ := json.Marshal(dto)
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateBook(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var errorResp handlers.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &errorResp)
		assert.NoError(t, err)
		assert.Equal(t, "Validation failed", errorResp.Error)

		mockUseCase.AssertExpectations(t)
	})
}

func TestBookHandler_UpdateBook(t *testing.T) {
	mockUseCase, handler := setupBookHandler()

	t.Run("successful update", func(t *testing.T) {
		bookID := uuid.New()
		dto := entities.UpdateBookDTO{
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

		mockUseCase.On("UpdateBook", mock.Anything, bookID, &dto).Return(expectedBook, nil).Once()

		reqBody, _ := json.Marshal(dto)
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/books/"+bookID.String(), bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(bookID.String())

		err := handler.UpdateBook(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var result entities.Book
		err = json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, expectedBook.Title, result.Title)

		mockUseCase.AssertExpectations(t)
	})
}

func TestBookHandler_DeleteBook(t *testing.T) {
	mockUseCase, handler := setupBookHandler()

	t.Run("successful deletion", func(t *testing.T) {
		bookID := uuid.New()

		mockUseCase.On("DeleteBook", mock.Anything, bookID).Return(nil).Once()

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/books/"+bookID.String(), nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(bookID.String())

		err := handler.DeleteBook(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var result handlers.SuccessResponse
		err = json.Unmarshal(rec.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, "Book deleted successfully", result.Message)

		mockUseCase.AssertExpectations(t)
	})
}