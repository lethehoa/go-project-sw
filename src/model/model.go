package model

import (
	// "gorm.io/gorm"
)

type User struct {
	Username string `gorm:"primaryKey;not null"`
	Email string `gorm:"not null"`
	Pub_key string `gorm:"not null"`
}