package models

import "time"

type User struct {
	ID        string     `gorm:"type:char(36);primaryKey"`
	Email     string     `gorm:"type:varchar(255);unique;not null"`
	Name      string     `gorm:"type:varchar(255);not null"`
	Phone     string     `gorm:"type:varchar(15)"`
	Role      string     `gorm:"type:enum('sender','member','both');default:'member'"`
	IsActive  bool       `gorm:"default:true"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"index"`
}
