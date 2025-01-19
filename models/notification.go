package models

import "time"

type Notification struct {
	ID          string    `gorm:"type:char(36);primaryKey"`
	SenderID    string    `gorm:"type:char(36);not null"`
	Subject     string    `gorm:"type:varchar(255);not null"`
	Message     string    `gorm:"type:text;not null"` // เปลี่ยนจาก JSON เป็น TEXT
	Priority    string    `gorm:"type:enum('low','normal','high');default:'normal'"`
	Status      string    `gorm:"type:enum('pending','in_progress','sent','failed');default:'pending'"`
	ScheduledAt time.Time `gorm:"type:datetime"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
