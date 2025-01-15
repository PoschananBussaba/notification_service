package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID        string         `gorm:"type:char(36);primary_key" json:"id"`
	Subject   string         `gorm:"not null" json:"subject"`
	Message   string         `gorm:"type:json;not null" json:"message"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
