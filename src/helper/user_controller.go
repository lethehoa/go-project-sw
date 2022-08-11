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
	var public_key_filepath string
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
		public_key_filepath = services.Create_ssh_key_file(body.Pub_key, body.Username)
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Create user successful"})
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
		Username string
		Email    string
		Pub_key  string
		Pri_key  string
	}

	c.BindJSON(&body)

	var user model.User
	DB.First(&user, "username = ?", username)

	DB.Model(&user).Updates(model.User{
		Username: body.Username,
		Email:    body.Email,
		Pub_key:  body.Pub_key,
	})
	c.IndentedJSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context) {
	username := c.Param("username")
	DB.Delete(&model.User{}, "username = ?", username)
	c.Status(200)
}
