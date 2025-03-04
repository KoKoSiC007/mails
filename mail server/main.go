package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"strings"

	"os"

	"github.com/emersion/go-msgauth/dkim"
	"github.com/emersion/go-smtp"
)

var domains = map[string]map[string]string{
	"ya.ru": {
		"addr": "localhost",
		"port": ":2525",
	},
	"mail.ru": {
		"addr": "localhost",
		"port": ":2626",
	},
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("Domain name must be present in args")
	}
	domain := os.Args[1]
	if domains[domain] == nil {
		log.Fatal("Unknown domain name")
	}

	s := smtp.NewServer(&Backend{Domain: domain})
	s.Addr = domains[domain]["port"]
	s.Domain = domains[domain]["addr"]
	//s.WriteTimeout = 10 * time.Second
	//s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	log.Println("Starting server at: ", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

type Backend struct {
	Domain string
}

func (bkd *Backend) NewSession(conn *smtp.Conn) (smtp.Session, error) {
	return &Session{HostDomain: bkd.Domain}, nil
}

type Session struct {
	From       string
	To         []string
	HostDomain string
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	fmt.Println("Mail from: ", from)
	s.From = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	fmt.Println("Rctp to: ", to)
	s.To = append(s.To, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if data, err := io.ReadAll(r); err != nil {
		return err
	} else {
		fmt.Println("Received message: ", string(data))

		for _, recipient := range s.To {
			toDomain := strings.Split(recipient, "@")[1]

			if s.HostDomain == toDomain {
				continue
			}
			if err := sendMail(s.From, recipient, data); err != nil {
				fmt.Printf("Failed to send email to %s: %v", recipient, err)
			} else {
				fmt.Printf("Email sent successfully to %s", recipient)
			}
		}
		return nil
	}
}

func (s *Session) AuthPlain(username, pass string) error {
	if username != "test" || pass != "test" {
		return fmt.Errorf("Invalid username or password")
	}

	return nil
}

func (s *Session) Logout() error {
	return nil
}

func (s *Session) Reset() {}

var dkimPrivateKey *rsa.PrivateKey

func init() {
	privateKeyPEM, err := os.ReadFile("keys")
	if err != nil {
		log.Fatalf("Failed to read private key: %v", err)
	}

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		log.Fatalf("Failed to parse PEM block containing the private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	dkimPrivateKey = privateKey
}

var dkimOptions = &dkim.SignOptions{
	Domain:   "example.com",
	Selector: "default",
	Signer:   dkimPrivateKey,
}

func sendMail(from string, to string, data []byte) error {
	domain := strings.Split(to, "@")[1]

	host := domains[domain]["addr"]
	port := domains[domain]["port"]
	address := fmt.Sprintf("%s%s", host, port)
	fmt.Println(address)
	c, err := smtp.Dial(address)
	if err != nil {
		return err
	}
	if err = c.Mail(from, nil); err != nil {
		c.Close()

		return err
	}
	if err = c.Rcpt(to, nil); err != nil {
		c.Close()
		return err
	}
	w, err := c.Data()
	if err != nil {
		c.Close()
		return err
	}
	_, err = w.Write(data)
	if err != nil {
		c.Close()
		return err
	}
	err = w.Close()
	if err != nil {
		c.Close()
		return err
	}

	c.Quit()

	return nil
}
