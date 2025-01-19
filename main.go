package main

import (
	"log"
	"notification_service/database"
	"notification_service/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// โหลด environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// ตั้งค่า Fiber App
	app := fiber.New()

	// เชื่อมต่อกับ Database
	database.ConnectDB()

	// ตั้งค่า Routes
	routes.SetupRoutes(app)

	// เริ่ม Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
