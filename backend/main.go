package main

import (
	"log"

	"byfood-library/internal/config"
	"byfood-library/internal/delivery/http/handlers"
	"byfood-library/internal/infrastructure/database"
	"byfood-library/internal/repositories"
	"byfood-library/internal/routes"
	"byfood-library/internal/usecases"
	_ "byfood-library/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// @title Book Library API
// @version 1.0
// @description A simple book library management API with UUID support and Echo framework
// @host localhost:8080
// @BasePath /
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize logger
	var logger *zap.Logger
	if cfg.Logging.Environment == "production" {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	// Initialize database connection
	db, err := database.InitDBWithConfig(cfg.Database.GetConnectionString())
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize Clean Architecture layers
	bookRepo := repositories.NewPostgresBookRepository(db)
	bookUseCase := usecases.NewBookUseCase(bookRepo, logger)
	bookHandler := handlers.NewBookHandler(bookUseCase, logger)
	urlHandler := handlers.NewURLHandler(logger)

	// Initialize Echo server
	e := echo.New()

	// Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.CORS.AllowedOrigins,
		AllowMethods: cfg.CORS.AllowedMethods,
		AllowHeaders: cfg.CORS.AllowedHeaders,
	}))

	// Setup routes with handlers
	handlers := &routes.Handlers{
		BookHandler: bookHandler,
		URLHandler:  urlHandler,
	}
	routes.SetupRoutes(e, cfg, handlers)

	// Start server
	address := cfg.Server.Host + ":" + cfg.Server.Port
	logger.Info("Server starting",
		zap.String("address", address),
		zap.String("environment", cfg.Logging.Environment),
		zap.String("framework", "echo"))

	if err := e.Start(address); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}