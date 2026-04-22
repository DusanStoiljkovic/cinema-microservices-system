package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	godotenv.Load()

	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbTcp := os.Getenv("DB_TCP")

	gormDB, err := gorm.Open(mysql.Open(dbUser+":"+dbPassword+dbTcp+dbName+"?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		fmt.Println("gorm DB connection ", err)
		return nil, err
	}

	return gormDB, nil
}
