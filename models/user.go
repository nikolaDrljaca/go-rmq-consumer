package models

import "time"

type User struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Email     string
	Messages []Message `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}