package main

import (
	"src/helper"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/create_user", helper.UserCreate)
	router.GET("/users", helper.ReturnAllUsers)
	router.GET("/user/:username", helper.ReturnSingleUser)
	router.PUT("/user/:username", helper.UserUpdate)
	router.DELETE("/user/:username", helper.DeleteUser)
	router.Run(":1234")
}
