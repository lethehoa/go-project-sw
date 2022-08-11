package services

import (
	"os"
	"log"
	"golang.org/x/crypto/ssh"
)

func Create_ssh_key_file(pub_key string, filename string) string{
	filename = filename + ".pub"
	f, _ := os.Create("/home/hoalt/script/ssh-key/" + filename)
	defer f.Close()

	_, err := f.WriteString(pub_key)

	if err != nil {
		log.Fatal("Failed to create session %v: ", err)
	}
	return filename
}

func Push_ssh_file_toSW(filepath string) {
	conn := establish_ssh_con()
	defer conn.Close()
	session, err := conn.NewSession()

	if err != nil {
		log.Fatal("Failed to create session %v: ", err)
	}

	// sdtin, _ := session.StdinPipe()
	session.Stdout = os.Stdout
	session.Stdin = os.Stdin
	session.Shell()
	// cmds := []string{
		// "ip a",
		// "cd /etc/sysconfig/network-script",
		// "pwd",
		// "ls -la",
	// }
	// for _, cmd := range cmds {
	// 	fmt.Fprintf(sdtin, "%s\n", cmd)
	// }

	session.Wait()
	session.Close()


}

func establish_ssh_con() *ssh.Client {
	config := &ssh.ClientConfig{
		User: "admin",
		Auth: []ssh.AuthMethod{
			ssh.Password("abc123"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", "192.168.2.41"+":22", config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	return conn
} 