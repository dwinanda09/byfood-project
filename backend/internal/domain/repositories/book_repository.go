package repositories

import (
	"context"

	"byfood-library/internal/domain/entities"
	"github.com/google/uuid"
)

type BookRepository interface {
	Create(ctx context.Context, book *entities.Book) (*entities.Book, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Book, error)
	GetAll(ctx context.Context) ([]*entities.Book, error)
	Update(ctx context.Context, id uuid.UUID, book *entities.Book) (*entities.Book, error)
	Delete(ctx context.Context, id uuid.UUID) error
}