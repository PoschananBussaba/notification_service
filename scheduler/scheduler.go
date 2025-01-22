package scheduler

import (
	"bytes"
	"encoding/base64"
	"fmt"
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

	"io/ioutil"

	"github.com/go-co-op/gocron"
)

func StartScheduler() {
	s := gocron.NewScheduler(time.UTC)

	// รัน Scheduler ทุก 1 นาที
	s.Every(1).Minute().Do(checkAndSendNotifications)

	// เริ่ม Scheduler
	s.StartAsync()
}

func checkAndSendNotifications() {
	var notifications []models.Notification

	// ดึง Notifications ที่มีสถานะ pending และ scheduled_at = เวลาปัจจุบัน
	currentTime := time.Now().UTC().Truncate(time.Minute) // ใช้ UTC และตัดเศษวินาทีออก
	if err := database.DB.Where("status = ? AND scheduled_at = ?", "pending", currentTime).Find(&notifications).Error; err != nil {
		log.Println("Failed to fetch notifications:", err)
		return
	}

	if len(notifications) == 0 {
		log.Println("No pending notifications to process.")
		return
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
			"status":     "sent",
			"updated_at": time.Now().UTC(),
		})
		log.Printf("Notification ID: %s sent successfully", notification.ID)
	}
}

func sendEmail(notification models.Notification) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// ดึงอีเมลของสมาชิกในกลุ่ม
	var emails []string
	query := `
		SELECT u.email
		FROM users u
		INNER JOIN target_group_members tgm ON u.id = tgm.user_id
		INNER JOIN target_groups tg ON tgm.group_id = tg.id
		WHERE tg.name = ? AND u.is_active = 1 AND tgm.is_active = 1
	`
	if err := database.DB.Raw(query, notification.TargetGroupName).Scan(&emails).Error; err != nil {
		log.Printf("Failed to fetch emails for target group %s: %v", notification.TargetGroupName, err)
		return err
	}

	if len(emails) == 0 {
		log.Printf("No active members in group: %s", notification.TargetGroupName)
		return nil
	}

	// สร้าง MIME message
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Header สำหรับ MIME
	headers := map[string]string{
		"From":         from,
		"To":           strings.Join(emails, ","),
		"Subject":      notification.Subject,
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
	textPart.Write([]byte(notification.Message))

	// แนบไฟล์
	attachments := strings.Split(notification.Attachments, ",")
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
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, emails, body.Bytes())
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Email sent successfully to: %v", emails)
	return nil
}
