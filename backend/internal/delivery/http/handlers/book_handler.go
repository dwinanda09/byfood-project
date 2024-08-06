package handlers

import (
	"net/http"

	"byfood-library/internal/domain/entities"
	"byfood-library/internal/middleware"
	"byfood-library/internal/usecases"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type bookHandler struct {
	bookUseCase usecases.BookUseCase
	logger      *zap.Logger
}

func NewBookHandler(bookUseCase usecases.BookUseCase, logger *zap.Logger) BookHandlerInterface {
	return &bookHandler{
		bookUseCase: bookUseCase,
		logger:      logger,
	}
}

// @Summary Get all books
// @Description Get all books from the library with UUID and timestamps
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {array} entities.Book
// @Failure 500 {object} ErrorResponse
// @Router /books [get]
func (h *bookHandler) GetBooks(c echo.Context) error {
	ctx := c.Request().Context()
	requestID := middleware.GetRequestID(c)
	ctx = middleware.WithRequestID(ctx, requestID)
	
	h.logger.Info("Getting all books", zap.String("request_id", requestID))
	
	books, err := h.bookUseCase.GetAllBooks(ctx)
	if err != nil {
		h.logger.Error("Failed to get all books", zap.String("request_id", requestID), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve books",
			Message: err.Error(),
		})
	}

	h.logger.Info("Successfully retrieved all books", zap.String("request_id", requestID), zap.Int("count", len(books)))
	return c.JSON(http.StatusOK, books)
}

// @Summary Get book by ID
// @Description Get a single book by its UUID
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book UUID"
// @Success 200 {object} entities.Book
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /books/{id} [get]
func (h *bookHandler) GetBook(c echo.Context) error {
	ctx := c.Request().Context()
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Error("Invalid UUID format", zap.String("id", idParam), zap.Error(err))
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid UUID format",
			Message: err.Error(),
		})
	}

	book, err := h.bookUseCase.GetBookByID(ctx, id)
	if err != nil {
		if err == entities.ErrBookNotFound {
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Book not found",
				Message: err.Error(),
			})
		}
		h.logger.Error("Failed to get book", zap.String("id", id.String()), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve book",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, book)
}

// @Summary Create a new book
// @Description Create a new book with title, author, and year
// @Tags books
// @Accept json
// @Produce json
// @Param book body entities.CreateBookDTO true "Book to create"
// @Success 201 {object} entities.Book
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /books [post]
func (h *bookHandler) CreateBook(c echo.Context) error {
	ctx := c.Request().Context()
	var dto entities.CreateBookDTO

	if err := c.Bind(&dto); err != nil {
		h.logger.Error("Failed to bind request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	book, err := h.bookUseCase.CreateBook(ctx, &dto)
	if err != nil {
		// Check for validation errors
		if err == entities.ErrInvalidTitle || err == entities.ErrInvalidAuthor || err == entities.ErrInvalidYear {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Validation failed",
				Message: err.Error(),
			})
		}
		h.logger.Error("Failed to create book", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create book",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, book)
}

// @Summary Update a book
// @Description Update an existing book by UUID
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book UUID"
// @Param book body entities.UpdateBookDTO true "Book data to update"
// @Success 200 {object} entities.Book
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /books/{id} [put]
func (h *bookHandler) UpdateBook(c echo.Context) error {
	ctx := c.Request().Context()
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Error("Invalid UUID format", zap.String("id", idParam), zap.Error(err))
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid UUID format",
			Message: err.Error(),
		})
	}

	var dto entities.UpdateBookDTO
	if err := c.Bind(&dto); err != nil {
		h.logger.Error("Failed to bind request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
	}

	book, err := h.bookUseCase.UpdateBook(ctx, id, &dto)
	if err != nil {
		if err == entities.ErrBookNotFound {
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Book not found",
				Message: err.Error(),
			})
		}
		// Check for validation errors
		if err == entities.ErrInvalidTitle || err == entities.ErrInvalidAuthor || err == entities.ErrInvalidYear {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "Validation failed",
				Message: err.Error(),
			})
		}
		h.logger.Error("Failed to update book", zap.String("id", id.String()), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update book",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, book)
}

// @Summary Delete a book
// @Description Delete a book by UUID
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book UUID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /books/{id} [delete]
func (h *bookHandler) DeleteBook(c echo.Context) error {
	ctx := c.Request().Context()
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Error("Invalid UUID format", zap.String("id", idParam), zap.Error(err))
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid UUID format",
			Message: err.Error(),
		})
	}

	err = h.bookUseCase.DeleteBook(ctx, id)
	if err != nil {
		if err == entities.ErrBookNotFound {
			return c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Book not found",
				Message: err.Error(),
			})
		}
		h.logger.Error("Failed to delete book", zap.String("id", id.String()), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to delete book",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Book deleted successfully",
	})
}