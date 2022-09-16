package util

import (
	"log"
	"net/http"
	"src/model"

	"src/helper"
	"src/services"

	"github.com/gin-gonic/gin"
)

var AllUsersForSync []model.User

func UploadFile(c *gin.Context) {
	file, _ := c.FormFile("file")
	path := "/home/hoalt/script/test-fol/upload/key.txt"
	c.SaveUploadedFile(file, path)
	c.String(http.StatusOK, "Files uploaded!")
}

func UpdateUser(c *gin.Context) {
	var deleteUser = SyncUser()
	log.Println(deleteUser)
	for _, i := range deleteUser {
		helper.DB.Delete(&model.User{}, "username = ?", i)
		services.Delete_sshkey_from_switch(i)
	}
	for _, i := range AllUser {
		if checkContainModelUser(AllUsersForSync, i.Username) {
			var user model.User
			helper.DB.First(&user, "username = ?", i.Username)
			helper.DB.Model(&user).Updates(model.User{
				Username: user.Username,
				Type:     user.Type,
				Email:    user.Username + "@vietnix.com.vn",
				Pub_key:  i.Key,
			})
			services.Create_ssh_key_file(i.Key, user.Username, 0)
		} else {
			addUser(i)
			services.Create_ssh_key_file(i.Key, i.Username, 1)
		}
	}
}

func SyncUser() []string {
	var deleteUser []string
	ParseUserKey()
	helper.DB.Find(&AllUsersForSync)
	for _, v := range AllUsersForSync {

		if !checkContain(AllUser, v.Username) {
			deleteUser = append(deleteUser, v.Username)
		}
	}
	return deleteUser
}

func checkContain(arr []User, element string) bool {
	for _, v := range arr {
		if v.Username == element {
			return true
		}
	}
	return false
}

func checkContainModelUser(arr []model.User, element string) bool {
	for _, v := range arr {
		if v.Username == element {
			return true
		}
	}
	return false
}

func addUser(userFromOutSide User) {
	email := userFromOutSide.Username + "@vietnix.com.vn"
	user := model.User{Username: userFromOutSide.Username, Type: userFromOutSide.Type, Email: email, Pub_key: userFromOutSide.Key}
	helper.DB.Create(&user)
}
