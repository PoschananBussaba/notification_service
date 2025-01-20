package models

import "time"

type User struct {
	ID        string     `gorm:"column:id;type:char(36);primaryKey"`
	Email     string     `gorm:"column:email;type:varchar(255);unique;not null"`
	Name      string     `gorm:"column:name;type:varchar(255);not null"`
	Phone     string     `gorm:"column:phone;type:varchar(15)"`
	Role      string     `gorm:"column:role;type:enum('sender','member','both');default:'member'"`
	IsActive  bool       `gorm:"column:is_active;default:true"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index"`
}
