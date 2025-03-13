package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"os"

	"github.com/emersion/go-msgauth/dkim"
	"github.com/emersion/go-smtp"
	senderservice "github.com/kokos/go-smtp-server/sender_service"
)

func main() {
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		log.Fatal("Unknown domain name")
	}
	address, err := GetMailAddress(domain)
	if err != nil {
		log.Fatal(err)
	}

	senderService := senderservice.NewSenderService()
	s := smtp.NewServer(&Backend{Domain: domain, SenderService: senderService})
	s.Addr = ":2525"
	s.Domain = address
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true
	log.Println("Starting server at: ", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

type Backend struct {
	Domain        string
	SenderService *senderservice.SenderService
}

func (bkd *Backend) NewSession(conn *smtp.Conn) (smtp.Session, error) {
	return &Session{HostDomain: bkd.Domain, senderService: bkd.SenderService}, nil
}

type Session struct {
	From          string
	To            []string
	HostDomain    string
	senderService *senderservice.SenderService
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

		s.processMessage(data)

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

func (s *Session) processMessage(data []byte) {
	for _, recipient := range s.To {
		toDomain := strings.Split(recipient, "@")[1]

		if s.HostDomain == toDomain {
			if err := s.verify(data); err != nil {
				fmt.Printf("Failed to verify mail, err: %v", err)
			} else {
				fmt.Printf("Mail verified!!!")
			}
			continue
		}

		if err := s.sendMail(s.From, recipient, data); err != nil {
			fmt.Printf("Failed to send email to %s: %v", recipient, err)
		} else {
			fmt.Printf("Email sent successfully to %s", recipient)
		}
	}
}

func (s *Session) verify(message []byte) error {
	r := strings.NewReader(string(message))

	verifications, err := dkim.Verify(r)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range verifications {
		if v.Err == nil {
			log.Println("Valid signature for:", v.Domain)
		} else {
			log.Println("Invalid signature for:", v.Domain, v.Err)
			return v.Err
		}
	}

	return nil
}

func (s *Session) sendMail(from string, to string, data []byte) error {
	domain := strings.Split(to, "@")[1]

	host, err := GetMailAddress(domain)
	if err != nil {
		return err
	}

	address := fmt.Sprintf("%s:2525", host)
	err = s.senderService.SendMail(address, from, to, data)
	if err != nil {
		return err
	}

	return nil
}

func GetMailAddress(domain string) (string, error) {
	r, err := net.LookupMX(domain)
	if err != nil {
		return "", err
	}
	if len(r) != 1 {
		return "", fmt.Errorf("DNS returns many MX records!")
	}
	return r[0].Host + domain, nil
}
