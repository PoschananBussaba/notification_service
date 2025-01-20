package main

import (
	"log"
	"notification_service/database"
	"notification_service/routes"
	"notification_service/scheduler"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	// เชื่อมต่อฐานข้อมูล
	database.ConnectDB()

	// เริ่ม Scheduler
	scheduler.StartScheduler()

	// ตั้งค่า Routes
	routes.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on http://localhost:" + port)
	log.Fatal(app.Listen(":" + port))
}
