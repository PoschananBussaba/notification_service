package scheduler

import (
	"log"
	"notification_service/database"
	"notification_service/models"
	"time"

	"github.com/go-co-op/gocron"
)

func StartScheduler() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Minute().Do(func() {
		var notifications []models.Notification
		now := time.Now()

		database.DB.Where("status = ? AND scheduled_at <= ?", "pending", now).Find(&notifications)

		for _, notification := range notifications {
			// Mock email sending
			log.Printf("Sending email for notification ID: %s\n", notification.ID)

			// Update status
			database.DB.Model(&notification).Update("status", "sent")
		}
	})

	s.StartAsync()
}
