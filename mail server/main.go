package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"os"

	"github.com/emersion/go-msgauth/dkim"
	"github.com/emersion/go-smtp"
	senderservice "github.com/kokos/go-smtp-server/sender_service"
)

var domains = map[string]map[string]string{
	"ya.ru": {
		"addr": "mail.ya.ru",
		"port": "2525",
	},
	"mail.ru": {
		"addr": "mail.mail.ru",
		"port": "2626",
	},
}

func main() {
	domain := os.Getenv("DOMAIN")
	if domains[domain] == nil {
		log.Fatal("Unknown domain name")
	}
	senderService := senderservice.NewSenderService()
	s := smtp.NewServer(&Backend{Domain: domain, SenderService: senderService})
	s.Addr = ":" + domains[domain]["port"]
	s.Domain = domains[domain]["addr"]
	//s.WriteTimeout = 10 * time.Second
	//s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true
	r, _ := net.LookupTXT("mail._domainkey.orion.ru.")
	fmt.Println(r)
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
	fmt.Println(string(message))
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

// func (s *Session) LookupTXT(domain string) ([]string, error) {
// 	var s1 string
// 	s1 = "brisbane._domainkey.mail.ru TXT \"v=DKIM1; k=rsa; t=s; p=LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUEwV3JpQzh2ejBzNGtyTWEvZnhVTwpCczh3a2l6OG9laGNoTnlxQTNOVHgrOVJMa1VTL1NraXdRdkJvOG50NERETmwvQWIrUmxaYm5CV0ZoYS9lZVNzCitjUmk3UWRMNmV5MWo2dklYamhVMVNJOTBFaEttWUtITzNDRUVoWVRrRkZPVmdlNVZqQXI2R0NoWTNCSGQ5QkcKUFRicmduN29kM05ad1h0RWM2N2J4MWRJdExkNnFqYzRxYW1aSmw0bXBkU0Nwcy92SVNuL2x5QUt5VFRKdEdRQwoyNEhjVlphY3VtRDZHZEJVZ1V3TkVOVVBRS0lCRlNiT1dta3hOMkU5UnJ6RkZQUUJIOGN6SkNrSVpYQjhTVGdwCnJ6QTBJNzV6enVibURZeGpRRElNbk50cHBoYkVpZ21ZczJ3VituVy8rOXJyTWZ0Q1BXR3RoMUg3RkxDWXlNRTUKeHdJREFRQUIKLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0tCg==\""
// 	return []string{s1}, nil
// }

func (s *Session) sendMail(from string, to string, data []byte) error {
	domain := strings.Split(to, "@")[1]

	host := domains[domain]["addr"]
	port, _ := strconv.Atoi(domains[domain]["port"])
	address := fmt.Sprintf("%s:%v", host, port)
	err := s.senderService.SendMail(address, from, to, data)
	if err != nil {
		return err
	}

	return nil
}
