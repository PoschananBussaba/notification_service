package models

import "time"

type TargetGroupMember struct {
	ID        string     `gorm:"column:id;type:char(36);primaryKey"`
	GroupID   string     `gorm:"column:group_id;type:char(36);not null"`
	UserID    string     `gorm:"column:user_id;type:char(36);not null"`
	Role      string     `gorm:"column:role;type:enum('member','admin');default:'member'"`
	IsActive  bool       `gorm:"column:is_active;default:true"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index"`
}
