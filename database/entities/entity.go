package entities

import (
	"fmt"

	"gorm.io/gorm"
)

var DB *gorm.DB

type BaseEntity interface {}

func MigrateEntities() error {
	fmt.Println("Migration started")
	err := DB.AutoMigrate(Author{}, Book{})
	
	return err
}