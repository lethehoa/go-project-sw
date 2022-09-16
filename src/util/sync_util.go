package util

import (
	"bufio"
	// "encoding/json"

	"log"
	"os"
	"regexp"
	"strings"
)

type User struct {
	Username string
	Type     string
	Key      string
}

// Constant Variable

var Administrator []User

var Tech_lv1 []User

var Tech_lv2 []User

var AllUser []User

var User_roles = []string{"Administrator", "Tech_lv2", "Tech_lv1"}

func ParseUserKey() {
	Reset()
	file := openFile("/home/hoalt/script/test-fol/upload/key.txt")
	scanner := bufio.NewScanner(file)

	for i := range User_roles {
		return_user_key_role(scanner, User_roles[i])
	}
	AllUser = append(AllUser, Administrator...)
	AllUser = append(AllUser, Tech_lv1...)
	AllUser = append(AllUser, Tech_lv2...)
	file.Close()
}

func Reset() {
	var temp []User

	Administrator = temp
	Tech_lv1 = temp
	Tech_lv2 = temp
	AllUser = temp
}

func return_user_key_role(scn *bufio.Scanner, role string) {
	var user_temp User
	for scn.Scan() {
		if strings.Contains(scn.Text(), "# == Start of "+role) {
			continue
		} else if strings.Contains(scn.Text(), "# == End of "+role) {
			break
		} else {
			text, check := return_user(scn.Text())
			if check {
				user_temp = User{text, role, ""}
			} else {
				user_temp.Key = text
				if strings.Compare(role, "Administrator") == 0 {
					Administrator = append(Administrator, user_temp)
				} else if strings.Compare(role, "Tech_lv2") == 0 {
					Tech_lv2 = append(Tech_lv2, user_temp)
				} else if strings.Compare(role, "Tech_lv1") == 0 {
					Tech_lv1 = append(Tech_lv1, user_temp)
				}

			}
		}
	}
}

func openFile(file_path string) *os.File {
	file, err := os.Open(file_path)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func return_user(str string) (string, bool) {
	rex, _ := regexp.MatchString("[^# ].*", str)
	if strings.Contains(str, "ssh-rsa") {
		rex = false
	}
	r, _ := regexp.Compile("[^# ].*")
	return r.FindString(str), rex
}
