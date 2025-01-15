package routes

import (
	"notification_service/handlers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Post("/notifications", handlers.CreateNotification)
}
