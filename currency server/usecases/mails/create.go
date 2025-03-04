package mails

import (
	"odyssey/m/v2/models/mail"
)

func (u *MailUseCase) Create(mail *mail.Mail, userId uint) (*mail.Mail, error) {
	user, err := u.userRepo.GetById(userId)
	if err != nil {
		return nil, err
	}
	mail.From = user.Email

	result, err := u.mailRepo.Create(mail, userId)
	if err != nil {
		return nil, err
	}

	mail.Id = result.ID
	err = u.mailService.Send(mail)
	if err != nil {
		return nil, err
	}
	return mail, nil
}
