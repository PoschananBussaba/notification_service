package scheduler

import (
	"log"
	"net/smtp"
	"os"
	"time"

	"notification_service/database"
	"notification_service/models"

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

	// ดึง Notifications ที่มีสถานะ pending และ scheduled_at <= เวลาปัจจุบัน
	if err := database.DB.Where("status = ? AND scheduled_at <= ?", "pending", time.Now()).Find(&notifications).Error; err != nil {
		log.Println("Failed to fetch notifications:", err)
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
		database.DB.Model(&notification).Update("status", "sent")
		log.Printf("Notification ID: %s sent successfully", notification.ID)
	}
}

func sendEmail(notification models.Notification) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	to := "linesunnyname@gmail.com" // ระบุผู้รับ (ดึงจากฐานข้อมูลหรือกำหนดคงที่)
	subject := "Subject: " + notification.Subject + "\n"
	message := "Message: " + notification.Message + "\n"
	body := []byte(subject + "\n" + message)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	// ส่งอีเมล
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, body)
}
