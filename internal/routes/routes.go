package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gunjan/user-api/internal/handler"
	"github.com/gunjan/user-api/internal/middleware"
	"go.uber.org/zap"
)

func SetupRoutes(app *fiber.App, userHandler *handler.UserHandler, logger *zap.Logger) {
	// Apply global middleware
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger(logger))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// User routes
	api := app.Group("/users")
	api.Post("/", userHandler.CreateUser)
	api.Get("/", userHandler.GetAllUsers)
	api.Get("/:id", userHandler.GetUserByID)
	api.Put("/:id", userHandler.UpdateUser)
	api.Delete("/:id", userHandler.DeleteUser)
}

