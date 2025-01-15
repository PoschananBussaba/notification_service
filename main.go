package main

import (
	"log"

	"notification_service/database"
	"notification_service/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// โหลด Environment Variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// เชื่อมต่อฐานข้อมูล
	database.ConnectDatabase()

	// สร้าง Fiber App
	app := fiber.New()

	// ตั้งค่าเส้นทาง API
	routes.Setup(app)

	// รันเซิร์ฟเวอร์
	log.Fatal(app.Listen(":8080"))
}
