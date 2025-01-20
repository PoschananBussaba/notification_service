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

	// สร้างเนื้อหาอีเมล
	subject := "Subject: " + notification.Subject + "\n"
	message := notification.Message + "\n"
	body := []byte(subject + "\n" + message)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	// ส่งอีเมลไปยังผู้รับทั้งหมด
	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, emails, body)
}
