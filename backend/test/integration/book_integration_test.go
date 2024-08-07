package integration

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"byfood-library/internal/config"
	"byfood-library/internal/domain/entities"
	"byfood-library/internal/repositories"
	"byfood-library/internal/usecases"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type BookIntegrationTestSuite struct {
	suite.Suite
	db        *sqlx.DB
	bookUC    usecases.BookUseCase
	bookRepo  repositories.BookRepository
	logger    *zap.Logger
}

func (s *BookIntegrationTestSuite) SetupSuite() {
	// Load test configuration
	cfg, err := config.Load("../../config/config.test.yaml")
	if err != nil {
		log.Fatal("Failed to load test configuration:", err)
	}

	// Setup logger
	s.logger, _ = zap.NewDevelopment()

	// Connect to test database
	s.db, err = sqlx.Connect("postgres", cfg.Database.URL)
	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}

	// Run migrations
	s.runMigrations()

	// Initialize repositories and use cases
	s.bookRepo = repositories.NewPostgresBookRepository(s.db, s.logger)
	s.bookUC = usecases.NewBookUseCase(s.bookRepo, s.logger)
}

func (s *BookIntegrationTestSuite) TearDownSuite() {
	s.db.Close()
}

func (s *BookIntegrationTestSuite) SetupTest() {
	// Clean database before each test
	s.db.Exec("DELETE FROM books")
}

func (s *BookIntegrationTestSuite) runMigrations() {
	migrations := []string{
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,
		`CREATE TABLE IF NOT EXISTS books (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			title VARCHAR(255) NOT NULL,
			author VARCHAR(255) NOT NULL,
			year INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
	}

	for _, migration := range migrations {
		_, err := s.db.Exec(migration)
		if err != nil {
			log.Fatal("Migration failed:", err)
		}
	}
}

func (s *BookIntegrationTestSuite) TestCreateBook() {
	dto := &entities.CreateBookDTO{
		Title:  "Integration Test Book",
		Author: "Test Author",
		Year:   2023,
	}

	book, err := s.bookUC.CreateBook(context.Background(), dto)

	s.NoError(err)
	s.NotNil(book)
	s.NotEmpty(book.ID)
	s.Equal(dto.Title, book.Title)
	s.Equal(dto.Author, book.Author)
	s.Equal(dto.Year, book.Year)
	s.False(book.CreatedAt.IsZero())
	s.False(book.UpdatedAt.IsZero())

	// Verify data in database
	var count int
	err = s.db.Get(&count, "SELECT COUNT(*) FROM books WHERE id = $1", book.ID)
	s.NoError(err)
	s.Equal(1, count)
}

func (s *BookIntegrationTestSuite) TestGetBookByID() {
	// Create a book first
	dto := &entities.CreateBookDTO{
		Title:  "Get Test Book",
		Author: "Get Test Author",
		Year:   2023,
	}

	createdBook, err := s.bookUC.CreateBook(context.Background(), dto)
	s.NoError(err)

	// Now get it by ID
	retrievedBook, err := s.bookUC.GetBookByID(context.Background(), createdBook.ID)

	s.NoError(err)
	s.NotNil(retrievedBook)
	s.Equal(createdBook.ID, retrievedBook.ID)
	s.Equal(createdBook.Title, retrievedBook.Title)
	s.Equal(createdBook.Author, retrievedBook.Author)
	s.Equal(createdBook.Year, retrievedBook.Year)
}

func (s *BookIntegrationTestSuite) TestGetBookByID_NotFound() {
	nonExistentID := uuid.New()

	book, err := s.bookUC.GetBookByID(context.Background(), nonExistentID)

	s.Error(err)
	s.Nil(book)
	s.Equal(entities.ErrBookNotFound, err)
}

func (s *BookIntegrationTestSuite) TestGetAllBooks() {
	// Create multiple books
	books := []*entities.CreateBookDTO{
		{Title: "Book 1", Author: "Author 1", Year: 2020},
		{Title: "Book 2", Author: "Author 2", Year: 2021},
		{Title: "Book 3", Author: "Author 3", Year: 2022},
	}

	var createdBooks []*entities.Book
	for _, dto := range books {
		book, err := s.bookUC.CreateBook(context.Background(), dto)
		s.NoError(err)
		createdBooks = append(createdBooks, book)
	}

	// Get all books
	allBooks, err := s.bookUC.GetAllBooks(context.Background())

	s.NoError(err)
	s.NotNil(allBooks)
	s.Len(allBooks, 3)

	// Verify books are ordered by created_at DESC
	s.True(allBooks[0].CreatedAt.After(allBooks[1].CreatedAt) || allBooks[0].CreatedAt.Equal(allBooks[1].CreatedAt))
	s.True(allBooks[1].CreatedAt.After(allBooks[2].CreatedAt) || allBooks[1].CreatedAt.Equal(allBooks[2].CreatedAt))
}

func (s *BookIntegrationTestSuite) TestUpdateBook() {
	// Create a book first
	dto := &entities.CreateBookDTO{
		Title:  "Original Book",
		Author: "Original Author",
		Year:   2020,
	}

	createdBook, err := s.bookUC.CreateBook(context.Background(), dto)
	s.NoError(err)

	// Wait a bit to ensure updated_at is different
	time.Sleep(10 * time.Millisecond)

	// Update the book
	updateDTO := &entities.UpdateBookDTO{
		Title:  "Updated Book",
		Author: "Updated Author",
		Year:   2023,
	}

	updatedBook, err := s.bookUC.UpdateBook(context.Background(), createdBook.ID, updateDTO)

	s.NoError(err)
	s.NotNil(updatedBook)
	s.Equal(createdBook.ID, updatedBook.ID)
	s.Equal(updateDTO.Title, updatedBook.Title)
	s.Equal(updateDTO.Author, updatedBook.Author)
	s.Equal(updateDTO.Year, updatedBook.Year)
	s.Equal(createdBook.CreatedAt.Unix(), updatedBook.CreatedAt.Unix())
	s.True(updatedBook.UpdatedAt.After(createdBook.UpdatedAt))
}

func (s *BookIntegrationTestSuite) TestUpdateBook_NotFound() {
	nonExistentID := uuid.New()
	updateDTO := &entities.UpdateBookDTO{
		Title:  "Updated Book",
		Author: "Updated Author",
		Year:   2023,
	}

	updatedBook, err := s.bookUC.UpdateBook(context.Background(), nonExistentID, updateDTO)

	s.Error(err)
	s.Nil(updatedBook)
	s.Equal(entities.ErrBookNotFound, err)
}

func (s *BookIntegrationTestSuite) TestDeleteBook() {
	// Create a book first
	dto := &entities.CreateBookDTO{
		Title:  "Book to Delete",
		Author: "Delete Author",
		Year:   2023,
	}

	createdBook, err := s.bookUC.CreateBook(context.Background(), dto)
	s.NoError(err)

	// Delete the book
	err = s.bookUC.DeleteBook(context.Background(), createdBook.ID)

	s.NoError(err)

	// Verify it's deleted
	_, err = s.bookUC.GetBookByID(context.Background(), createdBook.ID)
	s.Error(err)
	s.Equal(entities.ErrBookNotFound, err)

	// Verify data is removed from database
	var count int
	err = s.db.Get(&count, "SELECT COUNT(*) FROM books WHERE id = $1", createdBook.ID)
	s.NoError(err)
	s.Equal(0, count)
}

func (s *BookIntegrationTestSuite) TestDeleteBook_NotFound() {
	nonExistentID := uuid.New()

	err := s.bookUC.DeleteBook(context.Background(), nonExistentID)

	s.Error(err)
	s.Equal(entities.ErrBookNotFound, err)
}

func TestBookIntegrationTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	// Only run if TEST_DATABASE_URL is set
	if os.Getenv("TEST_DATABASE_URL") == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping integration tests")
	}

	suite.Run(t, new(BookIntegrationTestSuite))
}