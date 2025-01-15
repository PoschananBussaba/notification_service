package models

import (
	"time"

	"gorm.io/gorm"
)

type TargetGroupMember struct {
	ID        string         `gorm:"type:char(36);primary_key" json:"id"`
	GroupID   string         `gorm:"type:char(36);not null" json:"group_id"`
	UserID    string         `gorm:"type:char(36);not null" json:"user_id"`
	Role      string         `gorm:"type:enum('member', 'admin');default:'member'" json:"role"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
