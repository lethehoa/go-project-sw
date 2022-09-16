package main

import (
	"src/helper"
	"src/util"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 15

	router.POST("/create_user", helper.UserCreate)
	router.POST("/upload_file", util.UploadFile)

	router.GET("/users", helper.ReturnAllUsers)
	router.GET("/user/:username", helper.ReturnSingleUser)
	router.PUT("/user/:username", helper.UserUpdate)
	router.DELETE("/user/:username", helper.DeleteUser)

	router.GET("/user/test_getall", util.UpdateUser)
	router.Run(":1234")

}
