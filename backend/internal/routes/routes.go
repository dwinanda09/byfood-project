package routes

import (
	"byfood-library/internal/config"
	"byfood-library/internal/delivery/http/handlers"
	"byfood-library/internal/middleware"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Handlers struct {
	BookHandler handlers.BookHandlerInterface
	URLHandler  handlers.URLHandlerInterface
}

func SetupRoutes(e *echo.Echo, cfg *config.Config, h *Handlers) {
	// Request ID middleware applied first
	e.Use(middleware.DefaultMiddleware())

	// API version group - new versioned endpoints
	v1 := e.Group("/api/v1")
	booksGroup := v1.Group("/books")
	booksGroup.GET("", h.BookHandler.GetBooks)
	booksGroup.POST("", h.BookHandler.CreateBook)
	booksGroup.GET("/:id", h.BookHandler.GetBook)
	booksGroup.PUT("/:id", h.BookHandler.UpdateBook)
	booksGroup.DELETE("/:id", h.BookHandler.DeleteBook)

	// Legacy routes for backward compatibility with existing frontend
	e.GET("/books", h.BookHandler.GetBooks)
	e.POST("/books", h.BookHandler.CreateBook)
	e.GET("/books/:id", h.BookHandler.GetBook)
	e.PUT("/books/:id", h.BookHandler.UpdateBook)
	e.DELETE("/books/:id", h.BookHandler.DeleteBook)

	// Swagger documentation with configurable paths
	if cfg.API.EnableSwagger {
		e.GET(cfg.API.SwaggerPath+"/*", echoSwagger.WrapHandler)
		// Backward compatibility redirect
		e.GET("/docs", func(c echo.Context) error {
			return c.Redirect(302, cfg.API.SwaggerPath+"/")
		})
	}

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "OK"})
	})
}