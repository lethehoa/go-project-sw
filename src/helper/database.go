package helper

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Test123 = "Thehoa"

func ConnectToDB() {
	var err error
	username := os.Getenv("USER_DB")
	password := os.Getenv("PASSWORD")
	dsn := username + ":" + password + "@tcp(103.200.22.104:3306)/testDB"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Fail to connect to DB")
	}
}
