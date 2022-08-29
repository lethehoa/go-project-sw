package main

import (
	"src/helper"
	"src/model"
)

func init() {
	helper.LoadEnv()
	helper.ConnectToDB()
}

func main() {
	helper.DB.AutoMigrate(&model.User{})
}