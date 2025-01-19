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

	// ตรวจสอบ ScheduledAt
	if notification.ScheduledAt.IsZero() {
		notification.ScheduledAt = time.Now() // ตั้งค่าเริ่มต้นเป็นเวลาปัจจุบัน
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

	// ดึง Notifications ที่มีสถานะ pending และ scheduled_at <= เวลาปัจจุบัน
	if err := database.DB.Where("status = ? AND scheduled_at <= ?", "pending", time.Now()).Find(&notifications).Error; err != nil {
		log.Println("Failed to fetch notifications:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch notifications"})
	}

	// ส่งอีเมลสำหรับแต่ละ Notification
	for _, notification := range notifications {
		err := sendEmail(notification)
		if err != nil {
			log.Printf("Failed to send notification ID: %s, Error: %v", notification.ID, err)
			database.DB.Model(&notification).Update("status", "failed")
			continue
		}

		// อัปเดตสถานะเป็น sent
		database.DB.Model(&notification).Updates(map[string]interface{}{
			"status": "sent",
		})
		log.Printf("Notification ID: %s sent successfully", notification.ID)
	}

	return c.Status(200).JSON(fiber.Map{"message": "Notifications processed successfully"})
}

func sendEmail(notification models.Notification) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// ดึงอีเมลผู้รับจากฐานข้อมูลที่เกี่ยวข้องกับ target group
	var recipientEmails []string
	err := database.DB.Raw(`
		SELECT u.email 
		FROM users u
		INNER JOIN target_group_members tgm ON u.id = tgm.user_id
		INNER JOIN target_groups tg ON tgm.group_id = tg.id
		WHERE tg.id = ? AND u.is_active = 1
	`, notification.SenderID).Scan(&recipientEmails).Error

	if err != nil {
		log.Printf("Failed to fetch recipient emails for notification ID: %s, Error: %v", notification.ID, err)
		return err
	}

	// ตรวจสอบว่ามีผู้รับหรือไม่
	if len(recipientEmails) == 0 {
		log.Printf("No recipients found for notification ID: %s", notification.ID)
		return nil
	}

	// สร้างข้อความอีเมล
	subject := "Subject: " + notification.Subject + "\n"
	message := "Message: " + notification.Message + "\n"
	body := []byte(subject + "\n" + message)

	// ส่งอีเมลให้กับผู้รับทุกคน
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, recipientEmails, body)

	if err != nil {
		log.Printf("Failed to send notification ID: %s, Error: %v", notification.ID, err)
		return err
	}

	log.Printf("Notification ID: %s sent to %d recipients successfully", notification.ID, len(recipientEmails))
	return nil
}
