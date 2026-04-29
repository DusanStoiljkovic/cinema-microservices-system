package db

import (
	"fmt"
	"log"
	"movie-service/internal/models"
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
			fmt.Println("Connected to Movie Database...")

			err = db.AutoMigrate(&models.Movie{}, &models.Genre{})
			if err != nil {
				log.Fatal("Migration failed in Movie Database: ", err)
			} else {
				fmt.Println("Migrated successfully in Movie Database...")
			}

			return db, nil
		}
		fmt.Println("Waiting for Movie Database...")
		time.Sleep(2 * time.Second)
	}
	return nil, err
}
