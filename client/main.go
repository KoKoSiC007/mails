package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"log"
)

func main() {
	CreateSign()
	CheckSign()
}

func CreateSign() {
	privateKey, err := ReadKey()

	message := "Hello, Golang!"
	hash := sha256.Sum256([]byte(message))

	signature, err := rsa.SignPSS(
		rand.Reader,
		privateKey,
		crypto.SHA256,
		hash[:],
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Подпись: %x\n", signature)
	SendRequest(message, signature)
}

func ReadKey() (*rsa.PrivateKey, error) {
	privateKeyPEM, err := os.ReadFile("rsa_private_dev.pem")
	if err != nil {
		log.Fatalf("Failed to read private key: %v", err)
	}

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		log.Fatalf("Failed to parse PEM block containing the private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	parsedKey := privateKey.(*rsa.PrivateKey)
	return parsedKey, nil
}

func SendRequest(message string, signature []byte) {
	urlP, err := url.Parse("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	urlP.Path += "/api/v1/verify"
	parameters := url.Values{}
	parameters.Add("data", message)
	parameters.Add("signature", string(signature))
	urlP.RawQuery = parameters.Encode()
	resp, err := http.Get(urlP.String())
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 200 {
		fmt.Println("Message verify!")
	} else {
		fmt.Println("Message not verify!")
	}
}

type SignedMessage struct {
	Message   string `json:"message"`
	Signature []byte `json:"signature"`
}

func CheckSign() {
	pubKey := GetPublicKey()
	signedMessage := GetTestMessage()

	hash := sha256.Sum256([]byte(signedMessage.Message))
	err := rsa.VerifyPSS(
		pubKey,
		crypto.SHA256,
		hash[:],
		signedMessage.Signature,
		nil,
	)
	if err != nil {
		fmt.Errorf("Invalid signature")
	}
	fmt.Println("Signature valid!")
}

func GetPublicKey() *rsa.PublicKey {
	urlP, err := url.Parse("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	urlP.Path += "/api/v1/public_key"
	resp, err := http.Get(urlP.String())

	pubKeyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	pubKeyBlock, _ := pem.Decode(pubKeyBytes)
	if pubKeyBlock == nil {
		log.Fatal(fmt.Errorf("Parsing key error"))
	}

	pubKey, err := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	rsaPubKey := pubKey.(*rsa.PublicKey)
	return rsaPubKey
}

func GetTestMessage() *SignedMessage {
	urlP, err := url.Parse("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	urlP.Path += "/api/v1/test_message"
	resp, err := http.Get(urlP.String())

	var signedMessage SignedMessage
	decoder := json.NewDecoder(resp.Body)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&signedMessage)
	if err != nil {
		log.Fatal(err)
	}

	return &signedMessage
}
