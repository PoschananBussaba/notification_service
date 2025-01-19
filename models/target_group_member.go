package models

import "time"

type TargetGroupMember struct {
	ID        string     `gorm:"type:char(36);primaryKey"`
	GroupID   string     `gorm:"char(36);not null"`
	UserID    string     `gorm:"char(36);not null"`
	Role      string     `gorm:"type:enum('member','admin');default:'member'"`
	IsActive  bool       `gorm:"default:true"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"index"`
}
