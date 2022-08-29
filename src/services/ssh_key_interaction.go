package services

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"
	// "bufio"
	// "strings"
)

type Switch_VN struct {
	IP       string
	Password string
}

var current_switch []Switch_VN

func init() {
	//Load switch file
	f, _ := os.Open("/root/Documents/switch-information/information.txt")
	readInformationFromFile(f)
	f.Close()
	//Append to the slide current_switch
}

func Create_ssh_key_file(pub_key string, username string, check int) {
	filepath := "~/Documents/key/" + username + ".pub"
	// filenamefull := filename + ".pub"
	f, _ := os.Create(filepath)
	defer f.Close()
	_, err := f.WriteString(pub_key)
	if err != nil {
		log.Fatal("Fail to create key")
	}
	for i := 0; i < len(current_switch); i++ {
		c := exec.Command("scp", filepath, "admin@"+current_switch[i].IP+":bootflash:/ssh-key") //Đẩy key lên switch.
		in, err := session.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}
		var out bytes.Buffer
		var stderr bytes.Buffer
		c.Stdout = &out
		c.Stderr = &stderr
		err2 := c.Run()
		if err2 != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			return
		}
		fmt.Println("Result: " + out.String())
	}

	// in, err := session.StdinPipe()
	// if err != nil {
	//     log.Fatal(err)
	// }

	// _, err = in.Write([]byte(conn.password + "\n"))
	//             if err != nil {
	//                 break
	//             }

	if check == 1 {
		for i := 0; i < len(current_switch); i++ {
			Interact_ssh_key_in_SW(establish_ssh_con(current_switch[i].IP, current_switch[i].Password), username, "create")
		}

	} else {
		for i := 0; i < len(current_switch); i++ {
			Interact_ssh_key_in_SW(establish_ssh_con(current_switch[i].IP, current_switch[i].Password), username, "update")
		}
	}

}

func Delete_sshkey_from_switch(username string) {
	for i := 0; i < len(current_switch); i++ {
		Interact_ssh_key_in_SW(establish_ssh_con(current_switch[i].IP, current_switch[i].Password), username, "delete")
	}
}

func Interact_ssh_key_in_SW(conn *ssh.Client, username string, interaction string) {
	defer conn.Close()
	var cmds []string
	session, err := conn.NewSession()

	if err != nil {
		log.Fatal("Failed to create session %v: ", err)
	}

	sdtin, _ := session.StdinPipe()
	session.Stdout = os.Stdout
	session.Stdin = os.Stdin
	session.Shell()
	if interaction == "create" {
		cmds = []string{
			"conf t",
			"username " + username + " role network-admin", //Tạo user để ssh
			"username " + username + " sshkey file ssh-key/" + username + ".pub",
			"copy run start",
		}
	} else if interaction == "delete" {
		Disconnect_a_session(conn, username) // disconnect all the session of each user
		cmds = []string{
			"conf t",
			"no username " + username, //Xoá user ra khỏi switch
			"move ssh-key/" + username + ".pub old-ssh/" + username + ".pub",
			"copy run start",
		}
	} else if interaction == "update" {
		Disconnect_a_session(conn, username)
		cmds = []string{
			"conf t",
			"username " + username + " sshkey file ssh-key/" + username + ".pub",
			"move ssh-key/" + username + ".pub old-ssh/" + username + ".pub",
			"copy run start",
		}
	}
	for _, cmd := range cmds {
		fmt.Fprintf(sdtin, "%s\n", cmd)
	}
	session.Close()
}

func Disconnect_a_session(conn *ssh.Client, username string) {
	session, err := conn.NewSession()
	if err != nil {
		log.Fatal("Failed to create session %v: ", err)
	}

	sdtin, _ := session.StdinPipe()
	session.Stdin = os.Stdin
	session.Shell()
	fmt.Fprintf(sdtin, "%s\n", "clear user "+username)
}

func establish_ssh_con(ip string, password string) *ssh.Client {
	config := &ssh.ClientConfig{
		User: "admin",
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", ip+":22", config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	return conn
}

func readInformationFromFile(file *os.File) {
	scanner := bufio.NewScanner(file)
	var store_tmp []string
	for scanner.Scan() {
		store_tmp = strings.Split(scanner.Text(), " ")
		current_switch = append(current_switch, Switch_VN{store_tmp[0], store_tmp[1]})
	}
}
