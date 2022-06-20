package models

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBInit(host string, user string, pass string, name string, port string) *gorm.DB {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v", host, user, pass, name, port)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	
	if err != nil {
		panic("Failed to connect to database.")
	}

	//IMPORTANT - Ordering is important, think of relationships
	database.AutoMigrate(&User{}, &Message{})

	return database
}