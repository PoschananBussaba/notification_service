package handlers

import (
	"fmt"
	"log"
	"net/smtp"
	"notification_service/database"
	"notification_service/models"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetUsersByRole(c *fiber.Ctx) error {
	log.Println("Fetching users with role...")

	role := c.Query("role", "sender,both") // ค่าเริ่มต้น
	log.Println("Role filter:", role)

	// ใช้ roles ในคำสั่ง Where
	var users []models.User
	if err := database.DB.Where("role IN ?", strings.Split(role, ",")).Find(&users).Error; err != nil {
		log.Println("Error fetching users:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users"})
	}

	log.Println("Users found:", users)
	emails := make([]string, len(users))
	for i, user := range users {
		emails[i] = user.Email
	}

	return c.JSON(emails)
}

func GetTargetGroups(c *fiber.Ctx) error {
	var groupNames []string

	// Query ดึงเฉพาะคอลัมน์ name จาก target_groups
	if err := database.DB.Model(&models.TargetGroup{}).Pluck("name", &groupNames).Error; err != nil {
		log.Println("Error fetching target group names:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch target group names"})
	}

	return c.JSON(groupNames)
}

func GetEmailsByTargetGroup(groupName string) ([]string, error) {
	var emails []string

	query := `
		SELECT u.email
		FROM users u
		INNER JOIN target_group_members tgm ON u.id = tgm.user_id
		INNER JOIN target_groups tg ON tgm.group_id = tg.id
		WHERE tg.name = ? AND u.is_active = 1
	`
	if err := database.DB.Raw(query, groupName).Scan(&emails).Error; err != nil {
		return nil, err
	}

	return emails, nil
}

func GetEmailsByGroupName(groupName string) ([]string, error) {
	var emails []string

	// Query ดึงอีเมลจาก target_groups และ users
	query := `
		SELECT u.email
		FROM users u
		INNER JOIN target_group_members tgm ON u.id = tgm.user_id
		INNER JOIN target_groups tg ON tgm.group_id = tg.id
		WHERE tg.name = ? AND u.is_active = 1 AND tgm.is_active = 1
	`

	if err := database.DB.Raw(query, groupName).Scan(&emails).Error; err != nil {
		log.Printf("Error fetching emails for group %s: %v", groupName, err)
		return nil, err
	}

	return emails, nil
}

func CreateNotification(c *fiber.Ctx) error {
	var input struct {
		SenderEmail     string `json:"sender_email"`
		Subject         string `json:"subject"`
		Message         string `json:"message"`
		Priority        string `json:"priority"`
		ScheduledAt     string `json:"scheduled_at"`
		TargetGroupName string `json:"target_group_name"`
	}

	// Parse JSON Input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// ค้นหา SenderID จาก SenderEmail
	var user models.User
	if err := database.DB.Where("email = ?", input.SenderEmail).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Sender not found"})
	}

	// ตรวจสอบ Target Group Name
	var targetGroup models.TargetGroup
	if err := database.DB.Where("name = ?", input.TargetGroupName).First(&targetGroup).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Target group not found"})
	}

	// แปลงเวลาจาก input เป็น `time.Time`
	scheduledAt, err := time.Parse("2006-01-02 15:04:05.000", input.ScheduledAt)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid scheduled_at format"})
	}

	// สร้าง Notification
	notification := models.Notification{
		ID:              uuid.New().String(),
		SenderID:        user.ID,
		SenderEmail:     input.SenderEmail,
		Subject:         input.Subject,
		Message:         input.Message,
		Priority:        input.Priority,
		Status:          "pending",
		ScheduledAt:     scheduledAt,
		TargetGroupName: input.TargetGroupName,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	// บันทึก Notification ลงฐานข้อมูล
	if err := database.DB.Create(&notification).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create notification"})
	}

	return c.Status(201).JSON(notification)
}

// Helper Function
func parseTime(datetime string) time.Time {
	t, _ := time.Parse("2006-01-02T15:04:05", datetime)
	return t
}

func SendNotifications(c *fiber.Ctx) error {
	notificationID := c.Query("notification_id")

	if notificationID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Notification ID is required"})
	}

	// ดึงข้อมูล Notification จากฐานข้อมูล
	var notification models.Notification
	if err := database.DB.Where("id = ?", notificationID).First(&notification).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Notification not found"})
	}

	// ส่งข้อมูลไปยังกลุ่ม
	err := SendNotificationToGroup(notification.TargetGroupName, notification.Subject, notification.Message)
	if err != nil {
		log.Printf("Error sending notification: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to send notification"})
	}

	return c.JSON(fiber.Map{"message": "Notification sent successfully"})
}

func SendEmail(subject, message string, recipients []string) error {
	// อ่านค่าจาก .env
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// เตรียมเนื้อหาอีเมล
	body := fmt.Sprintf("Subject: %s\n\n%s", subject, message)

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// ส่งอีเมล
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, recipients, []byte(body))
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Email sent successfully to: %v", recipients)
	return nil
}

func SendNotificationsByGroup(c *fiber.Ctx) error {
	groupName := c.Query("group_name")
	subject := c.Query("subject")
	message := c.Query("message")

	if groupName == "" || subject == "" || message == "" {
		return c.Status(400).JSON(fiber.Map{"error": "group_name, subject, and message are required"})
	}

	// ส่งค่า groupName, subject และ message
	err := SendNotificationToGroup(groupName, subject, message)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to send notification"})
	}

	return c.JSON(fiber.Map{"message": "Notification sent successfully"})
}

func SendNotificationToGroup(groupName, subject, message string) error {
	// ดึงอีเมลผู้รับจากกลุ่ม
	emails, err := GetEmailsByGroupName(groupName)
	if err != nil {
		return err
	}

	if len(emails) == 0 {
		log.Printf("No active members in group: %s", groupName)
		return nil
	}

	// ส่งอีเมล
	err = SendEmail(subject, message, emails)
	if err != nil {
		log.Printf("Failed to send emails to group: %s", groupName)
		return err
	}

	log.Printf("Emails sent successfully to group: %s", groupName)
	return nil
}
