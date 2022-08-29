package helper

import (
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
)

var path_dir = "../src"

func LoadEnv() {
	err := godotenv.Load(filepath.Join(path_dir, ".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
