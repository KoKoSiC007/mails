package main

import (
	"fmt"
	"log"
	"net/smtp"
)

func main() {
	c, err := smtp.Dial("localhost:2626")
	if err != nil {
		fmt.Println(1)
		log.Fatal(err)
	}

	if err := c.Mail("kokos@mail.ru"); err != nil {
		fmt.Println(2)
		log.Fatal(err)
	}

	if err := c.Rcpt("kokos@ya.ru"); err != nil {
		fmt.Println(3)
		log.Fatal(err)
	}

	wc, err := c.Data()
	if err != nil {
		fmt.Println(4)
		log.Fatal(err)
	}

	_, err = fmt.Fprintf(wc, "This is email body")
	if err != nil {
		fmt.Println(5)
		log.Fatal(err)
	}

	err = wc.Close()
	if err != nil {
		fmt.Println(6)
		log.Fatal(err)
	}

	err = c.Quit()
	if err != nil {
		fmt.Println(7)
		log.Fatal(err)
	}
	fmt.Println("Message sended!")
}
