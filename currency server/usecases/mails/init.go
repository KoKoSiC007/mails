package mails

import (
	"odyssey/m/v2/internal/mail"
	"odyssey/m/v2/repositories"
)

type MailUseCase struct {
	mailRepo    *repositories.MailRepo
	userRepo    *repositories.UserRepo
	mailService *mail.MailService
}

func NewMailCase(mailRepo *repositories.MailRepo, userRepo *repositories.UserRepo, mailService *mail.MailService) (*MailUseCase, error) {
	return &MailUseCase{mailRepo: mailRepo, userRepo: userRepo, mailService: mailService}, nil
}
