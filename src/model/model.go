package model

// "gorm.io/gorm"

type User struct {
	Username string `gorm:"primaryKey;not null"`
	Type     string `gorm:"not null"`
	Email    string `gorm:"not null"`
	Pub_key  string `gorm:"not null"`
}
