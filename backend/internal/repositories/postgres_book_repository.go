package repositories

import (
	"context"

	"byfood-library/internal/domain/entities"
	"byfood-library/internal/domain/repositories"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type postgresBookRepository struct {
	db *sqlx.DB
}

func NewPostgresBookRepository(db *sqlx.DB) repositories.BookRepository {
	return &postgresBookRepository{db: db}
}

// Create using named parameters and struct scanning
func (r *postgresBookRepository) Create(ctx context.Context, book *entities.Book) (*entities.Book, error) {
	query := `INSERT INTO books (title, author, year) VALUES (:title, :author, :year) 
              RETURNING id, title, author, year, created_at, updated_at`

	var createdBook entities.Book
	rows, err := r.db.NamedQueryContext(ctx, query, book)
	if err != nil {
		return nil, entities.ErrDatabaseError
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&createdBook)
		if err != nil {
			return nil, entities.ErrDatabaseError
		}
	}
	return &createdBook, nil
}

// GetByID using Get for single row retrieval
func (r *postgresBookRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Book, error) {
	query := `SELECT id, title, author, year, created_at, updated_at FROM books WHERE id = $1`

	var book entities.Book
	err := r.db.GetContext(ctx, &book, query, id)
	if err != nil {
		return nil, entities.ErrBookNotFound
	}
	return &book, nil
}

// GetAll using Select for automatic slice population
func (r *postgresBookRepository) GetAll(ctx context.Context) ([]*entities.Book, error) {
	query := `SELECT id, title, author, year, created_at, updated_at FROM books ORDER BY created_at DESC`

	var books []entities.Book
	err := r.db.SelectContext(ctx, &books, query)
	if err != nil {
		return nil, entities.ErrDatabaseError
	}

	// Convert []entities.Book to []*entities.Book
	var bookPointers []*entities.Book
	for i := range books {
		bookPointers = append(bookPointers, &books[i])
	}
	return bookPointers, nil
}

// Update using named parameters
func (r *postgresBookRepository) Update(ctx context.Context, id uuid.UUID, book *entities.Book) (*entities.Book, error) {
	query := `UPDATE books SET title = :title, author = :author, year = :year 
              WHERE id = :id RETURNING id, title, author, year, created_at, updated_at`

	// Set the ID for the named query
	book.ID = id

	var updatedBook entities.Book
	rows, err := r.db.NamedQueryContext(ctx, query, book)
	if err != nil {
		return nil, entities.ErrDatabaseError
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&updatedBook)
		if err != nil {
			return nil, entities.ErrDatabaseError
		}
		return &updatedBook, nil
	}
	return nil, entities.ErrBookNotFound
}

// Delete using Exec
func (r *postgresBookRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM books WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return entities.ErrDatabaseError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return entities.ErrDatabaseError
	}

	if rowsAffected == 0 {
		return entities.ErrBookNotFound
	}

	return nil
}