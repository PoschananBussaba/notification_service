package routes

import (
	"notification_service/database"
	"notification_service/models"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	// Route สำหรับ GET Users
	app.Get("/users", func(c *fiber.Ctx) error {
		var users []models.User
		database.DB.Find(&users)
		return c.JSON(users)
	})

	// Route สำหรับ POST Users
	app.Post("/users", func(c *fiber.Ctx) error {
		user := new(models.User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		// บันทึกผู้ใช้ลงฐานข้อมูล
		if result := database.DB.Create(&user); result.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": result.Error.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(user)
	})
}
