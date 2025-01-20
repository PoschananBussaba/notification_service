package models

import "time"

type TargetGroup struct {
	ID          string     `gorm:"column:id;type:char(36);primaryKey"`
	Name        string     `gorm:"column:name;type:varchar(255);unique;not null"`
	Description string     `gorm:"column:description;type:text"`
	CreatedAt   time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;index"`
}
