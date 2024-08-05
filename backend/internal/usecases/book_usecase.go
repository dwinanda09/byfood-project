package usecases

import (
	"context"

	"byfood-library/internal/domain/entities"
	"byfood-library/internal/domain/repositories"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type BookUseCase interface {
	CreateBook(ctx context.Context, dto *entities.CreateBookDTO) (*entities.Book, error)
	GetBookByID(ctx context.Context, id uuid.UUID) (*entities.Book, error)
	GetAllBooks(ctx context.Context) ([]*entities.Book, error)
	UpdateBook(ctx context.Context, id uuid.UUID, dto *entities.UpdateBookDTO) (*entities.Book, error)
	DeleteBook(ctx context.Context, id uuid.UUID) error
}

type bookUseCase struct {
	bookRepo repositories.BookRepository
	logger   *zap.Logger
}

func NewBookUseCase(bookRepo repositories.BookRepository, logger *zap.Logger) BookUseCase {
	return &bookUseCase{
		bookRepo: bookRepo,
		logger:   logger,
	}
}

func (uc *bookUseCase) CreateBook(ctx context.Context, dto *entities.CreateBookDTO) (*entities.Book, error) {
	// Validate DTO
	if err := dto.Validate(); err != nil {
		uc.logger.Error("Validation failed for CreateBookDTO", zap.Error(err))
		return nil, err
	}

	// Convert DTO to entity
	book := dto.ToBook()

	// Create book through repository
	createdBook, err := uc.bookRepo.Create(ctx, book)
	if err != nil {
		uc.logger.Error("Failed to create book", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("Book created successfully", zap.String("id", createdBook.ID.String()))
	return createdBook, nil
}

func (uc *bookUseCase) GetBookByID(ctx context.Context, id uuid.UUID) (*entities.Book, error) {
	book, err := uc.bookRepo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to get book by ID", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	return book, nil
}

func (uc *bookUseCase) GetAllBooks(ctx context.Context) ([]*entities.Book, error) {
	books, err := uc.bookRepo.GetAll(ctx)
	if err != nil {
		uc.logger.Error("Failed to get all books", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("Retrieved books successfully", zap.Int("count", len(books)))
	return books, nil
}

func (uc *bookUseCase) UpdateBook(ctx context.Context, id uuid.UUID, dto *entities.UpdateBookDTO) (*entities.Book, error) {
	// Validate DTO
	if err := dto.Validate(); err != nil {
		uc.logger.Error("Validation failed for UpdateBookDTO", zap.Error(err))
		return nil, err
	}

	// Create book entity from DTO
	book := &entities.Book{
		Title:  dto.Title,
		Author: dto.Author,
		Year:   dto.Year,
	}

	// Update book through repository
	updatedBook, err := uc.bookRepo.Update(ctx, id, book)
	if err != nil {
		uc.logger.Error("Failed to update book", zap.String("id", id.String()), zap.Error(err))
		return nil, err
	}

	uc.logger.Info("Book updated successfully", zap.String("id", updatedBook.ID.String()))
	return updatedBook, nil
}

func (uc *bookUseCase) DeleteBook(ctx context.Context, id uuid.UUID) error {
	err := uc.bookRepo.Delete(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to delete book", zap.String("id", id.String()), zap.Error(err))
		return err
	}

	uc.logger.Info("Book deleted successfully", zap.String("id", id.String()))
	return nil
}