package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"log"
)

func main() {
	// Парсим приватный ключ
	privateKey, err := ReadKey()

	// Сообщение для подписи
	message := "Hello, Golang!"

	// Вычисляем хеш сообщения
	hash := sha256.Sum256([]byte(message))

	// Создаем подпись
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
