package mail

import (
	"fmt"
	"net/smtp"
	"odyssey/m/v2/models/mail"
)

type Config struct {
	Addr string
}

type MailService struct {
	config Config
}

func NewMailService(config *Config) *MailService {
	return &MailService{config: *config}
}

func (s *MailService) Send(mail *mail.Mail) error {
	c, err := smtp.Dial(s.config.Addr)
	if err != nil {
		fmt.Println(1)
		return err
	}

	if err := c.Mail(mail.From); err != nil {
		fmt.Println(2)
		return err
	}

	if err := c.Rcpt(mail.To); err != nil {
		fmt.Println(3)
		return err
	}

	wc, err := c.Data()
	if err != nil {
		fmt.Println(4)
		return err
	}

	_, err = fmt.Fprintf(wc, mail.Body)
	if err != nil {
		fmt.Println(5)
		return err
	}

	err = wc.Close()
	if err != nil {
		fmt.Println(6)
		return err
	}

	err = c.Quit()
	if err != nil {
		fmt.Println(7)
		return err
	}
	fmt.Println("Message sended!")
	return nil
}
