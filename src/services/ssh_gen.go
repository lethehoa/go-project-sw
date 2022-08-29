package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	// "io/ioutil"
	// "log"
	// "os"
	"bufio"
	"bytes"
	"os"

	"golang.org/x/crypto/ssh"
)

type Switch_VN struct {
	IP       string
	Password string
}

var current_switch []Switch_VN

func Generate_ssh_key_pair(passphrase string) (public_key string, private_key string) {
	var PrivatekeyRow bytes.Buffer
	//Generate
	pri_key, _ := rsa.GenerateKey(rand.Reader, 512)
	privatekeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(pri_key),
	}

	privatekeyPEM, err := x509.EncryptPEMBlock(rand.Reader, privatekeyPEM.Type, privatekeyPEM.Bytes, []byte(passphrase), x509.PEMCipherAES256)
	//Convert privatekey to string
	err = pem.Encode(&PrivatekeyRow, privatekeyPEM)

	//Create public key
	pub, err := ssh.NewPublicKey(&pri_key.PublicKey)
	if err != nil {
		return PrivatekeyRow.String(), string(ssh.MarshalAuthorizedKey(pub))
	}

	// return PrivatekeyRow.String(), string(ssh.MarshalAuthorizedKey(pub))
	return "tess_public", "test_private"
}

func main() {
	// 	f, _ := os.Open("/home/hoalt/Documents/switch-information/information.txt")
	// 	readInformationFromFile(f)
	// 	fmt.Println(current_switch)
	// 	f.Close()
}

func readInformationFromFile(file *os.File) {
	scanner := bufio.NewScanner(file)
	var store_tmp []string
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		store_tmp = strings.Split(scanner.Text(), " ")

		current_switch = append(current_switch, Switch_VN{store_tmp[0], store_tmp[1]})
	}
}
