package services

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	// "io/ioutil"
	// "log"
	// "os"
	"bytes"
	"golang.org/x/crypto/ssh"
)


func Generate_ssh_key_pair(passphrase string) (public_key string, private_key string){
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

