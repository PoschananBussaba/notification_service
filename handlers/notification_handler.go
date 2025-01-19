package handlers

import (
	"log"
	"net/smtp"
	"notification_service/database"
	"notification_service/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// generateUUID - สร้าง UUID
func generateUUID() string {
	return uuid.New().String()
}

func CreateNotification(c *fiber.Ctx) error {
	var notification models.Notification

	// แปลงข้อมูลจาก JSON
	if err := c.BodyParser(&notification); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// ตรวจสอบและแปลงค่าของ `scheduled_at`
	if notification.ScheduledAt != nil {
		layout := "2006-01-02 15:04" // รูปแบบที่มาจากฟอร์ม HTML
		parsedTime, err := time.Parse(layout, notification.ScheduledAt.Format(layout))
		if err != nil {
			log.Println("Invalid datetime format:", err)
			return c.Status(400).JSON(fiber.Map{"error": "Invalid datetime format"})
		}
		notification.ScheduledAt = &parsedTime
	}

	// สร้าง UUID
	notification.ID = generateUUID()

	// บันทึกลงฐานข้อมูล
	if err := database.DB.Create(&notification).Error; err != nil {
		log.Println("Failed to create notification:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create notification"})
	}

	return c.Status(201).JSON(notification)
}

func SendNotifications(c *fiber.Ctx) error {
	var notifications []models.Notification

	// ดึง Notifications ที่สถานะเป็น pending
	if err := database.DB.Where("status = ?", "pending").Find(&notifications).Error; err != nil {
		log.Println("Failed to fetch notifications:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch notifications"})
	}

	// อ่านค่าจาก .env
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	for _, notification := range notifications {
		// เตรียมข้อความ
		to := "recipient@example.com" // ระบุอีเมลผู้รับ (แก้ไขตามต้องการ)
		subject := "Subject: " + notification.Subject + "\n"
		message := notification.Message + "\n"
		body := []byte(subject + "\n" + message)

		// ส่งอีเมล
		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, body)
		if err != nil {
			log.Printf("Failed to send notification ID: %s, Error: %v", notification.ID, err)
			database.DB.Model(&notification).Update("status", "failed")
			continue
		}

		// อัปเดตสถานะเป็น sent
		database.DB.Model(&notification).Updates(map[string]interface{}{
			"status":       "sent",
			"scheduled_at": time.Now(),
		})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Notifications processed successfully"})
}
