package services

import (
	"os"
	"log"
	"golang.org/x/crypto/ssh"
	"fmt"
	"os/exec"
)

func Create_ssh_key_file(pub_key string, filename string){ 			// thực hiện hai vấn đề, thứ nhất là tạo file ssh-key, thứ hai là copy sshkey vào sw
	filename = "D:\\Documents\\Learn_Go\\ssh-key\\" + filename + ".pub"
	f, _ := os.Create(filename)
	defer f.Close()
	_, err := f.WriteString(pub_key)
	if err != nil {
		log.Fatal("Failed to create session %v: ", err)
	}
	c := exec.Command("scp", "-i", "D:\\Documents\\Learn_Go\\ssh-key\\switch-test", filename ,"admin@192.168.57.150:bootflash:")
	if err := c.Run(); err != nil { 
        log.Fatal("Failed to create session %v: ", err)
    }
	Enable_ssh_key_in_SW(establish_ssh_con())
}

func Enable_ssh_key_in_SW(conn *ssh.Client) {
	defer conn.Close()
	session, err := conn.NewSession()

	if err != nil {
		log.Fatal("Failed to create session %v: ", err)
	}

	sdtin, _ := session.StdinPipe()
	session.Stdout = os.Stdout
	session.Stdin = os.Stdin
	session.Shell()
	cmds := []string{
		"conf t",
		"username admin sshkey file ta.nhattan.pub",
	}
	for _, cmd := range cmds {
		fmt.Fprintf(sdtin, "%s\n", cmd)
	}

	session.Wait()
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