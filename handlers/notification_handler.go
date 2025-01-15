package handlers

import (
	"encoding/json"
	"notification_service/database"
	"notification_service/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type NotificationRequest struct {
	SenderID     string   `json:"sender_id"`
	Subject      string   `json:"subject"`
	Message      string   `json:"message"`
	Priority     string   `json:"priority"`
	RecipientIDs []string `json:"recipient_ids"`
}

func CreateNotification(c *fiber.Ctx) error {
	var req NotificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// แปลง message ให้เป็น JSON String
	messageJSON, err := json.Marshal(map[string]string{"text": req.Message})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse message"})
	}

	notification := models.Notification{
		ID:        uuid.New().String(),
		SenderID:  req.SenderID,
		Subject:   req.Subject,
		Message:   string(messageJSON), // ใช้ JSON ที่แปลงแล้ว
		Priority:  req.Priority,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	// บันทึก Notification ในฐานข้อมูล
	if err := database.DB.Create(&notification).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create notification"})
	}

	// เพิ่มผู้รับ
	for _, userID := range req.RecipientIDs {
		recipient := models.NotificationRecipient{
			ID:             uuid.New().String(),
			NotificationID: notification.ID,
			UserID:         userID,
			Channel:        "email",
			DeliveryStatus: "pending",
			CreatedAt:      time.Now(),
		}
		database.DB.Create(&recipient)
	}

	return c.Status(fiber.StatusCreated).JSON(notification)
}
