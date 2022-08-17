package services

import (
	"os"
	"log"
	"golang.org/x/crypto/ssh"
	"fmt"
	"os/exec"
	"bytes"
	"bufio"
	"strings"
)

// user admin sẽ  hoàn toàn được ssh vào thông qua key mà không cần pass, hoặc sẽ chạy 1 lệnh ssh để copy vào.

func Create_ssh_key_file(pub_key string, username string){ 			// thực hiện hai vấn đề, thứ nhất là tạo file ssh-key, thứ hai là copy sshkey vào sw
	filepath := "D:\\Documents\\Learn_Go\\ssh-key\\" + username + ".pub"   // code trên Windows, sẽ fix lại sau khi test và chuyển sang linux
	// filenamefull := filename + ".pub"
	f, _ := os.Create(filepath)
	defer f.Close()
	_, err := f.WriteString(pub_key)
	if err != nil {
		log.Fatal("Fail to create key")
	}
	c := exec.Command("scp", "-i", "D:\\Documents\\Learn_Go\\ssh-key\\switch-nopass", filepath ,"admin@192.168.57.150:bootflash:/ssh-key")		//Đẩy key lên switch. // code trên Windows, sẽ fix lại sau khi test và chuyển sang linux
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
	Enable_ssh_key_in_SW(establish_ssh_con(), username, "create")
}

func Delete_sshkey_from_switch(username string) {
	Enable_ssh_key_in_SW(establish_ssh_con(), username, "delete")
}

func Enable_ssh_key_in_SW(conn *ssh.Client, username string ,interaction string) {
	defer conn.Close()
	var cmds []string
	var tmp string
	session, err := conn.NewSession()

	if err != nil {
		log.Fatal("Failed to create session %v: ", err)
	}

	sdtin, _ := session.StdinPipe()
	sdtout, _ := session.StdoutPipe()
	session.Stdout = os.Stdout
	session.Stdin = os.Stdin
	session.Shell()
	if interaction == "create" {
		cmds = []string{
			"conf t",
			"username " + username + " role network-admin",		//Tạo user để ssh
			"username " + username + " sshkey file ssh-key/" + username + ".pub",
			"copy run start",
		}
	} else {
		Disconnect_a_session(conn, username) // disconnect all the session of each user 
		cmds = []string{
			"conf t",
			"no username " + username,		//Xoá user ra khỏi switch
			"move ssh-key/" + username + ".pub old-ssh/" + username + ".pub",
			"copy run start",
		}
		scanner := bufio.NewScanner(sdtout)
		for scanner.Scan() {
			tmp = scanner.Text()
			if (strings.Contains(tmp, username)){
				// fmt.Println(text)
				x := strings.Split(tmp, " ")
				fmt.Println(len(x))
			}
		}
	}
	
	for _, cmd := range cmds {
		fmt.Fprintf(sdtin, "%s\n", cmd)
	}
	session.Close()
}

func Disconnect_a_session(conn *ssh.Client, username string){
	defer conn.Close()
	// var cmds []string
	session, err := conn.NewSession()
	if err != nil {
		log.Fatal("Failed to create session %v: ", err)
	}

	sdtin, _ := session.StdinPipe()
	session.Stdin = os.Stdin
	session.Shell()
	fmt.Fprintf(sdtin, "%s\n", "clear user " + username)
	session.Close()
}

func establish_ssh_con() *ssh.Client {
	config := &ssh.ClientConfig{
		User: "admin",
		Auth: []ssh.AuthMethod{
			ssh.Password("admin"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", "192.168.57.150"+":22", config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	return conn
} 