package models

import (
	"time"

	"gorm.io/gorm"
)

type TargetGroup struct {
	ID          string              `gorm:"type:char(36);primary_key" json:"id"`
	Name        string              `gorm:"type:varchar(255);not null" json:"name"`
	Description string              `gorm:"type:text" json:"description"`
	CreatedAt   time.Time           `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time           `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt      `gorm:"index" json:"-"`
	Members     []TargetGroupMember `gorm:"foreignKey:GroupID"`
}
