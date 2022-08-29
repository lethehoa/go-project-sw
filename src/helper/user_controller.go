package helper

import (
	"fmt"
	"net/http"
	"src/model"
	"github.com/gin-gonic/gin"
	"src/services"
)

var Users []model.User

func init() {
	LoadEnv()
	ConnectToDB()
}

func UserCreate(c *gin.Context) {
	// var public_key_filepath string
	var body struct {
		Username string
		Pub_key  string
	}
	c.BindJSON(&body) //Accept input from user

	email := body.Username + "@vietnix.com.vn"
	user := model.User{Username: body.Username, Email: email, Pub_key: body.Pub_key}
	result := DB.Create(&user)
	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Duplicate username, please try again",
		})
	} else {
		services.Create_ssh_key_file(body.Pub_key, body.Username, 1)
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Create and add user successful"})
	}
}

func ReturnAllUsers(c *gin.Context) {
	DB.Find(&Users)
	c.IndentedJSON(http.StatusOK, Users)

}

func ReturnSingleUser(c *gin.Context) {
	username := c.Param("username")
	DB.Find(&Users)
	var user model.User

	fmt.Println(username)
	result := DB.First(&user, "username = ?", username)
	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "User is not available"})
	} else {
		c.IndentedJSON(http.StatusOK, user)
	}
}

func UserUpdate(c *gin.Context) {
	username := c.Param("username")
	
	var body struct {
		Pub_key  string
	}

	c.BindJSON(&body)

	var user model.User
	result := DB.First(&user, "username = ?", username)
	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "User is not available"})
		return
	} else {
		DB.Model(&user).Updates(model.User{
			Username: username,
			Email: username + "@vietnix.com.vn",
			Pub_key:  body.Pub_key,
		})
		services.Create_ssh_key_file(body.Pub_key, username, 0)
		c.IndentedJSON(http.StatusOK, user)
	}
	
}

func DeleteUser(c *gin.Context) {
	username := c.Param("username")
	DB.Delete(&model.User{}, "username = ?", username)
	services.Delete_sshkey_from_switch(username)
	c.Status(200)
}
