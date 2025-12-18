package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gunjan/user-api/config"
	"github.com/gunjan/user-api/internal/handler"
	"github.com/gunjan/user-api/internal/logger"
	"github.com/gunjan/user-api/internal/repository"
	"github.com/gunjan/user-api/internal/routes"
	"github.com/gunjan/user-api/internal/service"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	if err := logger.InitLogger(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Log.Sync()

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db, err := sql.Open("mysql", cfg.Database.DSN())
	if err != nil {
		logger.Log.Fatal("Failed to open database connection", zap.Error(err))
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		logger.Log.Fatal("Failed to ping database", zap.Error(err))
	}
	logger.Log.Info("Database connection established")

	// Initialize repository
	userRepo := repository.NewUserRepository(db, logger.Log)

	// Initialize service
	userService := service.NewUserService(userRepo, logger.Log)

	// Initialize handler
	userHandler := handler.NewUserHandler(userService, logger.Log)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New())

	// Setup routes
	routes.SetupRoutes(app, userHandler, logger.Log)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	logger.Log.Info("Starting server", zap.String("address", addr))
	if err := app.Listen(addr); err != nil {
		logger.Log.Fatal("Failed to start server", zap.Error(err))
	}
}

