package services

import (
	"os"
	"log"
	"golang.org/x/crypto/ssh"
	"fmt"
	"os/exec"
	"bytes"
	// "bufio"
	// "strings"
)

func Create_ssh_key_file(pub_key string, username string, check int){ 		
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
	if check == 1 {
		Interact_ssh_key_in_SW(establish_ssh_con(), username, "create")
	} else {
		Interact_ssh_key_in_SW(establish_ssh_con(), username, "update")
	}
	
}

func Delete_sshkey_from_switch(username string) {
	Interact_ssh_key_in_SW(establish_ssh_con(), username, "delete")
}

func Interact_ssh_key_in_SW(conn *ssh.Client, username string ,interaction string) {
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
			"username " + username + " role network-admin",		//Tạo user để ssh
			"username " + username + " sshkey file ssh-key/" + username + ".pub",
			"copy run start",
		}
	} else if interaction == "delete" {
		Disconnect_a_session(conn, username) // disconnect all the session of each user 
		cmds = []string{
			"conf t",
			"no username " + username,		//Xoá user ra khỏi switch
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

func Disconnect_a_session(conn *ssh.Client, username string){
	session, err := conn.NewSession()
	if err != nil {
		log.Fatal("Failed to create session %v: ", err)
	}

	sdtin, _ := session.StdinPipe()
	session.Stdin = os.Stdin
	session.Shell()
	fmt.Fprintf(sdtin, "%s\n", "clear user " + username)
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