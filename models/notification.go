package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID          string                  `gorm:"type:char(36);primary_key" json:"id"`
	SenderID    string                  `gorm:"type:char(36);not null" json:"sender_id"`
	Subject     string                  `gorm:"type:varchar(255);not null" json:"subject"`
	Message     string                  `gorm:"type:json;not null" json:"message"`
	Status      string                  `gorm:"type:enum('pending', 'in_progress', 'sent', 'failed');default:'pending'" json:"status"`
	Priority    string                  `gorm:"type:enum('low', 'normal', 'high');default:'normal'" json:"priority"`
	ScheduledAt *time.Time              `json:"scheduled_at"`
	CreatedAt   time.Time               `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time               `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt          `gorm:"index" json:"-"`
	Recipients  []NotificationRecipient `gorm:"foreignKey:NotificationID"`
}

type NotificationRecipient struct {
	ID               string         `gorm:"type:char(36);primary_key" json:"id"`
	NotificationID   string         `gorm:"type:char(36);not null" json:"notification_id"`
	UserID           string         `gorm:"type:char(36);not null" json:"user_id"`
	Channel          string         `gorm:"type:enum('email', 'line', 'sms');not null" json:"channel"`
	DeliveryStatus   string         `gorm:"type:enum('pending', 'delivered', 'failed');default:'pending'" json:"delivery_status"`
	DeliveryAttempts int            `gorm:"default:0" json:"delivery_attempts"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}
