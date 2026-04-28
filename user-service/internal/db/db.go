package db

import (
	"fmt"
	"log"
	"os"
	"time"
	"user-service/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			fmt.Println("Connected to DB")

			err = db.AutoMigrate(&models.User{})
			if err != nil {
				log.Fatal("Migration failed", err)
			} else {
				fmt.Println("Migrated successfully...")
			}

			return db, nil
		}
		fmt.Println("Waiting for DB...")
		time.Sleep(2 * time.Second)
	}

	return nil, err
}
