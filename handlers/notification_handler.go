package handlers

import (
	"log"
	"notification_service/database"
	"notification_service/models"
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

	// แปลง ScheduledAt เป็นรูปแบบที่ถูกต้อง
	scheduledAt, err := time.Parse("2006-01-02T15:04", input.ScheduledAt)
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
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// บันทึก Notification ลงฐานข้อมูล
	if err := database.DB.Create(&notification).Error; err != nil {
		log.Println("Failed to create notification:", err)
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
	var notifications []models.Notification

	// ดึง Notifications ที่สถานะเป็น pending และ scheduled_at <= เวลาปัจจุบัน
	if err := database.DB.Where("status = ? AND scheduled_at <= ?", "pending", time.Now()).Find(&notifications).Error; err != nil {
		log.Println("Failed to fetch notifications:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch notifications"})
	}

	for _, notification := range notifications {
		// Logic สำหรับการส่ง notification (อาจเป็นการส่งอีเมลหรืออื่น ๆ)
		log.Printf("Sending notification ID: %s", notification.ID)

		// อัปเดตสถานะเป็น sent
		if err := database.DB.Model(&notification).Update("status", "sent").Error; err != nil {
			log.Printf("Failed to update status for notification ID: %s", notification.ID)
		}
	}

	return c.Status(200).JSON(fiber.Map{"message": "Notifications sent successfully"})
}
