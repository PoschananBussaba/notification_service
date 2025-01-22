package models

import "time"

type Notification struct {
	ID              string    `gorm:"column:id;type:char(36);primaryKey"`
	SenderID        string    `gorm:"column:sender_id;type:char(36);not null"`
	SenderEmail     string    `gorm:"column:sender_email;type:varchar(255)"`
	Subject         string    `gorm:"column:subject;type:varchar(255);not null"`
	Message         string    `gorm:"column:message;type:text;not null"`
	Priority        string    `gorm:"column:priority;type:enum('low','normal','high');default:'normal'"`
	Status          string    `gorm:"column:status;type:enum('pending','in_progress','sent','failed');default:'pending'"`
	ScheduledAt     time.Time `gorm:"column:scheduled_at;type:datetime;not null"`
	Attachments     string    `gorm:"column:attachments;type:text"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime"`
	TargetGroupName string    `gorm:"column:target_group_name;type:varchar(255)"`
}
