package tests

import (
	"context"
	"testing"

	"byfood-library/internal/domain/entities"
	domain_repositories "byfood-library/internal/domain/repositories"
	"byfood-library/internal/repositories"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setupRepositoryTest() (*sqlx.DB, sqlmock.Sqlmock, domain_repositories.BookRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	// Wrap with sqlx.NewDb for sqlx compatibility
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewPostgresBookRepository(sqlxDB)

	return sqlxDB, mock, repo
}

func TestPostgresBookRepository_Create(t *testing.T) {
	db, mock, repo := setupRepositoryTest()
	defer db.Close()

	t.Run("successful creation", func(t *testing.T) {
		book := &entities.Book{
			Title:  "The Go Programming Language",
			Author: "Alan Donovan",
			Year:   2015,
		}

		expectedID := uuid.New()
		rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "created_at", "updated_at"}).
			AddRow(expectedID, "The Go Programming Language", "Alan Donovan", 2015, "2024-01-01T00:00:00Z", "2024-01-01T00:00:00Z")

		mock.ExpectQuery(`INSERT INTO books \(title, author, year\) VALUES \(\$1, \$2, \$3\)`).
			WithArgs("The Go Programming Language", "Alan Donovan", 2015).
			WillReturnRows(rows)

		result, err := repo.Create(context.Background(), book)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedID, result.ID)
		assert.Equal(t, "The Go Programming Language", result.Title)
		assert.Equal(t, "Alan Donovan", result.Author)
		assert.Equal(t, 2015, result.Year)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		book := &entities.Book{
			Title:  "The Go Programming Language",
			Author: "Alan Donovan",
			Year:   2015,
		}

		mock.ExpectQuery(`INSERT INTO books \(title, author, year\) VALUES \(\$1, \$2, \$3\)`).
			WithArgs("The Go Programming Language", "Alan Donovan", 2015).
			WillReturnError(sqlmock.ErrCancelled)

		result, err := repo.Create(context.Background(), book)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, entities.ErrDatabaseError, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgresBookRepository_GetByID(t *testing.T) {
	db, mock, repo := setupRepositoryTest()
	defer db.Close()

	t.Run("successful retrieval", func(t *testing.T) {
		bookID := uuid.New()
		rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "created_at", "updated_at"}).
			AddRow(bookID, "The Go Programming Language", "Alan Donovan", 2015, "2024-01-01T00:00:00Z", "2024-01-01T00:00:00Z")

		mock.ExpectQuery(`SELECT id, title, author, year, created_at, updated_at FROM books WHERE id = \$1`).
			WithArgs(bookID).
			WillReturnRows(rows)

		result, err := repo.GetByID(context.Background(), bookID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, bookID, result.ID)
		assert.Equal(t, "The Go Programming Language", result.Title)
		assert.Equal(t, "Alan Donovan", result.Author)
		assert.Equal(t, 2015, result.Year)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("book not found", func(t *testing.T) {
		bookID := uuid.New()

		mock.ExpectQuery(`SELECT id, title, author, year, created_at, updated_at FROM books WHERE id = \$1`).
			WithArgs(bookID).
			WillReturnError(sqlmock.ErrCancelled)

		result, err := repo.GetByID(context.Background(), bookID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, entities.ErrBookNotFound, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgresBookRepository_GetAll(t *testing.T) {
	db, mock, repo := setupRepositoryTest()
	defer db.Close()

	t.Run("successful retrieval", func(t *testing.T) {
		bookID1 := uuid.New()
		bookID2 := uuid.New()
		rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "created_at", "updated_at"}).
			AddRow(bookID1, "Book 1", "Author 1", 2020, "2024-01-01T00:00:00Z", "2024-01-01T00:00:00Z").
			AddRow(bookID2, "Book 2", "Author 2", 2021, "2024-01-02T00:00:00Z", "2024-01-02T00:00:00Z")

		mock.ExpectQuery(`SELECT id, title, author, year, created_at, updated_at FROM books ORDER BY created_at DESC`).
			WillReturnRows(rows)

		result, err := repo.GetAll(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, bookID1, result[0].ID)
		assert.Equal(t, "Book 1", result[0].Title)
		assert.Equal(t, bookID2, result[1].ID)
		assert.Equal(t, "Book 2", result[1].Title)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, title, author, year, created_at, updated_at FROM books ORDER BY created_at DESC`).
			WillReturnError(sqlmock.ErrCancelled)

		result, err := repo.GetAll(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, entities.ErrDatabaseError, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgresBookRepository_Update(t *testing.T) {
	db, mock, repo := setupRepositoryTest()
	defer db.Close()

	t.Run("successful update", func(t *testing.T) {
		bookID := uuid.New()
		book := &entities.Book{
			Title:  "Updated Title",
			Author: "Updated Author",
			Year:   2022,
		}

		rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "created_at", "updated_at"}).
			AddRow(bookID, "Updated Title", "Updated Author", 2022, "2024-01-01T00:00:00Z", "2024-01-01T01:00:00Z")

		mock.ExpectQuery(`UPDATE books SET title = \$1, author = \$2, year = \$3 WHERE id = \$4`).
			WithArgs("Updated Title", "Updated Author", 2022, bookID).
			WillReturnRows(rows)

		result, err := repo.Update(context.Background(), bookID, book)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, bookID, result.ID)
		assert.Equal(t, "Updated Title", result.Title)
		assert.Equal(t, "Updated Author", result.Author)
		assert.Equal(t, 2022, result.Year)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("book not found", func(t *testing.T) {
		bookID := uuid.New()
		book := &entities.Book{
			Title:  "Updated Title",
			Author: "Updated Author",
			Year:   2022,
		}

		rows := sqlmock.NewRows([]string{"id", "title", "author", "year", "created_at", "updated_at"})

		mock.ExpectQuery(`UPDATE books SET title = \$1, author = \$2, year = \$3 WHERE id = \$4`).
			WithArgs("Updated Title", "Updated Author", 2022, bookID).
			WillReturnRows(rows)

		result, err := repo.Update(context.Background(), bookID, book)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, entities.ErrBookNotFound, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgresBookRepository_Delete(t *testing.T) {
	db, mock, repo := setupRepositoryTest()
	defer db.Close()

	t.Run("successful deletion", func(t *testing.T) {
		bookID := uuid.New()

		mock.ExpectExec(`DELETE FROM books WHERE id = \$1`).
			WithArgs(bookID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(context.Background(), bookID)

		assert.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("book not found", func(t *testing.T) {
		bookID := uuid.New()

		mock.ExpectExec(`DELETE FROM books WHERE id = \$1`).
			WithArgs(bookID).
			WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected

		err := repo.Delete(context.Background(), bookID)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrBookNotFound, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		bookID := uuid.New()

		mock.ExpectExec(`DELETE FROM books WHERE id = \$1`).
			WithArgs(bookID).
			WillReturnError(sqlmock.ErrCancelled)

		err := repo.Delete(context.Background(), bookID)

		assert.Error(t, err)
		assert.Equal(t, entities.ErrDatabaseError, err)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}