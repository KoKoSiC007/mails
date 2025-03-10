package senderservice

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	smtp_client "net/smtp"
	"os"

	"github.com/emersion/go-msgauth/dkim"
	"gopkg.in/gomail.v2"
)

type SenderService struct {
	dkimOptions *dkim.SignOptions
}

func NewSenderService() *SenderService {
	privateKeyPEM, err := os.ReadFile("keys")
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

	dkimOptions := &dkim.SignOptions{
		Domain:   "mail.ru",
		Selector: "mail",
		Signer:   parsedKey,
	}

	return &SenderService{dkimOptions: dkimOptions}
}

func (s *SenderService) SendMail(address, from, to string, data []byte) error {
	mail := s.buildMail(from, to, data)
	message, err := s.signMail(mail)
	if err != nil {
		return err
	}

	err = smtp_client.SendMail(address, nil, from, []string{to}, []byte(message))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *SenderService) buildMail(from, to string, body []byte) *gomail.Message {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, "")
	m.SetHeader("To", m.FormatAddress(to, ""), to)
	m.SetBody("text/txt", string(body))
	return m
}

func (s *SenderService) signMail(mail *gomail.Message) (string, error) {
	raw_mail := new(bytes.Buffer)
	mail.WriteTo(raw_mail)

	var raw_signed_mail bytes.Buffer

	if err := dkim.Sign(&raw_signed_mail, bytes.NewReader(raw_mail.Bytes()), s.dkimOptions); err != nil {
		return "", fmt.Errorf("Failed to sign email with DKIM: %v", err)
	}
	return raw_signed_mail.String(), nil
}
