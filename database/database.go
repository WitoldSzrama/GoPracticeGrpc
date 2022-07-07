package database

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"practice/three/database/entities"
)

var DB *gorm.DB

func OpenConnection() (db *gorm.DB) {
	DB, err := gorm.Open(mysql.New(mysql.Config{
		DSN: os.Getenv("DSN"),
	}))

	if err != nil {
		log.Panicf("Problem with database connection: %v\n", err)
	}

	entities.DB = DB
	return DB
}

func MigrateEntities() error {
	println("Migrate started")
	err := DB.AutoMigrate(&entities.Author{}, &entities.Book{})

	return err
}