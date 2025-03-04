package repositories

import (
	"database/sql"
	"odyssey/m/v2/models/mail"
	"odyssey/m/v2/models/user"

	"gorm.io/gorm"
)

type dbMail struct {
	gorm.Model
	ID     uint `gorm:"primaryKey"`
	UserId uint
	User   user.User `gorm:"references:Id"`
	To     string    `gorm:"not null"`
	Data   string    `gorm:"not null"`
}

type MailRepo struct {
	db *gorm.DB
}

func NewMailRepo(db *gorm.DB) (*MailRepo, error) {
	db.AutoMigrate(&dbMail{})

	return &MailRepo{db: db}, nil
}

func (repo *MailRepo) Create(mail *mail.Mail, userId uint) (*dbMail, error) {
	user := &user.User{Id: userId}

	entity := dbMail{
		To:   mail.To,
		Data: mail.Body,
		User: *user,
	}
	result := repo.db.Create(&entity)
	if result.Error != nil {
		return nil, result.Error
	}

	return &entity, nil
}

func (repo *MailRepo) GetByUserId(userId uint) (*[]dbMail, error) {
	var mails []dbMail

	result := repo.db.Where("user_id = @user_id", sql.Named("user_id", userId)).Find(&mails)
	if result.Error != nil {
		return nil, result.Error
	}

	return &mails, nil
}
