package routes

import (
	"notification_service/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Route หลัก
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./frontend/index.html") // เส้นทางไปยังไฟล์ HTML
	})

	// API Routes
	api := app.Group("/api")

	api.Post("/notifications", handlers.CreateNotification)
	api.Get("/send-notifications", handlers.SendNotifications)

}
