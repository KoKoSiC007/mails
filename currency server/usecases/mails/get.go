package mails

import (
	"odyssey/m/v2/models/mail"
)

func (u *MailUseCase) Get(userId uint) (*mail.MailsData, error) {
	result, err := u.mailRepo.GetByUserId(userId)
	if err != nil {
		return nil, err
	}

	var mails []mail.Mail
	for _, val := range *result {
		mail := mail.Mail{
			Id:   val.ID,
			To:   val.To,
			Body: val.Data,
		}

		mails = append(mails, mail)
	}

	return &mail.MailsData{Messages: mails}, nil
}
