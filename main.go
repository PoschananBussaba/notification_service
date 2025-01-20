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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	// เชื่อมต่อฐานข้อมูล
	database.ConnectDB()

	// ตั้งค่า Routes
	routes.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on http://localhost:" + port)
	log.Fatal(app.Listen(":" + port))
}

/*package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// โหลดไฟล์ .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	fmt.Println("SMTP_EMAIL:", from)
	fmt.Println("SMTP_PASSWORD:", password)
	fmt.Println("SMTP_HOST:", smtpHost)
	fmt.Println("SMTP_PORT:", smtpPort)

	if from == "" || password == "" || smtpHost == "" || smtpPort == "" {
		log.Fatal("One or more SMTP environment variables are missing.")
	}

	to := []string{"linesunnyname@gmail.com"}
	subject := "Subject: Test Email\n"
	body := "This is a test email."
	message := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	fmt.Println("Email sent successfully to:", to)
}*/
