package database

import (
	"log"
	"notification_service/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "root:password@tcp(127.0.0.1:3306)/notification_service?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// บันทึก Database instance
	DB = database

	// เรียก AutoMigrate เพื่อสร้างตาราง
	err = DB.AutoMigrate(
		&models.User{},
		&models.Notification{},
		&models.NotificationRecipient{},
		&models.TargetGroup{},
		&models.TargetGroupMember{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database connected and migrated successfully!")
}
