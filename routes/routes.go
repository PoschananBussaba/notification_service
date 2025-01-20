package routes

import (
	"notification_service/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Route หลัก
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./frontend/index.html")
	})

	// API Routes
	api := app.Group("/api")

	api.Post("/notifications", handlers.CreateNotification)
	api.Get("/send-notifications", handlers.SendNotifications)
	api.Get("/users", handlers.GetUsersByRole)
	api.Get("/target-groups", handlers.GetTargetGroups)
}
