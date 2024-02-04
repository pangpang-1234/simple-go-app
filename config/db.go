package config

import (
	"log"
	"os"

	// "simplegoapp/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDB() {
	var err error

	dsn := os.Getenv("DATABASE_CONNECTION")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	_ = db
}

func GetDB() *gorm.DB{
	return db
}

func CloseDB() {
	dbInstance, _ := db.DB()
    _ = dbInstance.Close()
}
