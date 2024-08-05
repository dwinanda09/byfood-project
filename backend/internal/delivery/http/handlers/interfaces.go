package handlers

import "github.com/labstack/echo/v4"

// BookHandlerInterface for testing
type BookHandlerInterface interface {
	GetBooks(c echo.Context) error
	GetBook(c echo.Context) error
	CreateBook(c echo.Context) error
	UpdateBook(c echo.Context) error
	DeleteBook(c echo.Context) error
}

// URLHandlerInterface for URL processing functionality
type URLHandlerInterface interface {
	ProcessURL(c echo.Context) error
}

// ErrorResponse for consistent error handling
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse for consistent success responses
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}