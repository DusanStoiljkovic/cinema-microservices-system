package db

import (
	"booking-service/internal/models"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			fmt.Println("Connected to Booking Database...")

			err = db.AutoMigrate(&models.Hall{}, &models.Projection{}, &models.Ticket{}, &models.Order{})
			if err != nil {
				log.Fatal("Migration failed in Booking Database...")
			} else {
				fmt.Println("Migrated successfully in Booking Database...")
			}

			return db, nil
		}
		fmt.Println("Waiting for Booking Databse...")
		time.Sleep(2 * time.Second)
	}
	return nil, err
}
