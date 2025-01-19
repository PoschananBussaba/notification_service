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

	// เชื่อมต่อกับ Database
	database.ConnectDB()

	// ตั้งค่า Routes
	routes.SetupRoutes(app)

	// เริ่ม Scheduler
	scheduler.StartScheduler()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
