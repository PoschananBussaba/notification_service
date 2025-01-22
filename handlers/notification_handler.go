package handlers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"notification_service/database"
	"notification_service/models"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetUsersByRole(c *fiber.Ctx) error {
	role := c.Query("role", "sender,both")

	var users []models.User
	if err := database.DB.Where("role IN ?", strings.Split(role, ",")).Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users"})
	}

	emails := make([]string, len(users))
	for i, user := range users {
		emails[i] = user.Email
	}

	return c.JSON(emails)
}

func GetTargetGroups(c *fiber.Ctx) error {
	var groupNames []string

	if err := database.DB.Model(&models.TargetGroup{}).Pluck("name", &groupNames).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch target group names"})
	}

	return c.JSON(groupNames)
}

func GetEmailsByGroupName(groupName string) ([]string, error) {
	var emails []string

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
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Failed to parse form data"})
	}

	// ดึงข้อมูลไฟล์
	files := form.File["attachments"] // ชื่อ input file ในฟอร์ม
	var filePaths []string

	for _, file := range files {
		filePath := "./uploads/" + file.Filename
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save file"})
		}
		filePaths = append(filePaths, filePath)
	}

	// รับข้อมูลฟอร์มอื่น
	input := struct {
		SenderEmail     string `form:"sender_email"`
		Subject         string `form:"subject"`
		Message         string `form:"message"`
		Priority        string `form:"priority"`
		ScheduledAt     string `form:"scheduled_at"`
		TargetGroupName string `form:"target_group_name"`
	}{}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse form"})
	}

	// ดึง Sender ID จาก Sender Email
	var user models.User
	if err := database.DB.Where("email = ?", input.SenderEmail).First(&user).Error; err != nil {
		log.Printf("Sender email not found: %s", input.SenderEmail)
		return c.Status(404).JSON(fiber.Map{"error": "Sender not found"})
	}

	// ตรวจสอบว่า Sender ID มีค่า
	if user.ID == "" {
		log.Println("Sender ID is empty for email:", input.SenderEmail)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid sender email"})
	}

	// แปลงเวลาจาก ISO8601
	scheduledAt, err := time.Parse(time.RFC3339, input.ScheduledAt)
	if err != nil {
		log.Printf("Error parsing scheduled_at: %v", input.ScheduledAt)
		return c.Status(400).JSON(fiber.Map{"error": "Invalid scheduled_at format"})
	}

	// สร้าง Notification
	notification := models.Notification{
		ID:              uuid.New().String(),
		SenderID:        user.ID, // เพิ่ม Sender ID ที่ถูกต้อง
		SenderEmail:     input.SenderEmail,
		Subject:         input.Subject,
		Message:         input.Message,
		Priority:        input.Priority,
		ScheduledAt:     scheduledAt,
		Attachments:     strings.Join(filePaths, ","), // เก็บ path ของไฟล์
		TargetGroupName: input.TargetGroupName,
		Status:          "pending",
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	if err := database.DB.Create(&notification).Error; err != nil {
		log.Printf("Database error: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create notification"})
	}

	return c.Status(201).JSON(notification)
}

func SendEmail(subject, message string, recipients []string, attachments []string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Header สำหรับ MIME
	headers := map[string]string{
		"From":         from,
		"To":           strings.Join(recipients, ","),
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": `multipart/mixed; boundary="` + writer.Boundary() + `"`,
	}

	for key, value := range headers {
		body.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	body.WriteString("\r\n")

	// ส่วนข้อความ
	textPart, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": {"text/plain; charset=\"utf-8\""},
	})
	if err != nil {
		log.Printf("Failed to create text part: %v", err)
		return err
	}
	textPart.Write([]byte(message))

	// แนบไฟล์
	for _, filePath := range attachments {
		fileName := filepath.Base(filePath)
		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to read file %s: %v", filePath, err)
			continue
		}

		log.Printf("Attaching file: %s", filePath)

		filePart, err := writer.CreatePart(textproto.MIMEHeader{
			"Content-Disposition":       {fmt.Sprintf("attachment; filename=\"%s\"", fileName)},
			"Content-Type":              {"application/octet-stream"},
			"Content-Transfer-Encoding": {"base64"},
		})
		if err != nil {
			log.Printf("Failed to create MIME part for file %s: %v", fileName, err)
			continue
		}

		// เขียนไฟล์ใน Base64
		encoded := base64.StdEncoding.EncodeToString(fileContent)
		filePart.Write([]byte(encoded))
	}

	writer.Close()

	// ส่งอีเมล
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, recipients, body.Bytes())
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Email sent successfully to: %v", recipients)
	return nil
}

func SendNotifications(c *fiber.Ctx) error {
	notificationID := c.Query("notification_id")

	if notificationID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Notification ID is required"})
	}

	var notification models.Notification
	if err := database.DB.Where("id = ?", notificationID).First(&notification).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Notification not found"})
	}

	err := SendNotificationToGroup(notification.TargetGroupName, notification.Subject, notification.Message, strings.Split(notification.Attachments, ","))
	if err != nil {
		log.Printf("Error sending notification: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to send notification"})
	}

	return c.JSON(fiber.Map{"message": "Notification sent successfully"})
}

func SendNotificationsByGroup(c *fiber.Ctx) error {
	groupName := c.Query("group_name")
	subject := c.Query("subject")
	message := c.Query("message")

	if groupName == "" || subject == "" || message == "" {
		return c.Status(400).JSON(fiber.Map{"error": "group_name, subject, and message are required"})
	}

	err := SendNotificationToGroup(groupName, subject, message, []string{})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to send notification"})
	}

	return c.JSON(fiber.Map{"message": "Notification sent successfully"})
}

func SendNotificationToGroup(groupName, subject, message string, attachments []string) error {
	emails, err := GetEmailsByGroupName(groupName)
	if err != nil {
		return err
	}

	if len(emails) == 0 {
		log.Printf("No active members in group: %s", groupName)
		return nil
	}

	// ส่งอีเมลพร้อมไฟล์แนบ
	err = SendEmail(subject, message, emails, attachments)
	if err != nil {
		log.Printf("Failed to send emails to group: %s", groupName)
		return err
	}

	log.Printf("Emails sent successfully to group: %s", groupName)
	return nil
}
