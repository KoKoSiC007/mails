package repositories

import (
	users "odyssey/m/v2/models/user"

	"gorm.io/gorm"
)

type dbUser struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Email     string `gorm:"not null"`
	Password  string `gorm:"not null"`
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) (*UserRepo, error) {
	db.AutoMigrate(&dbUser{})

	return &UserRepo{db: db}, nil
}

func (repo *UserRepo) Create(user *users.NewUser) (*dbUser, error) {
	entity := dbUser{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	}
	result := repo.db.Create(&entity)
	if result.Error != nil {
		return nil, result.Error
	}

	return &entity, nil
}

func (repo *UserRepo) Get(login, pass string) (*dbUser, error) {
	var user dbUser
	result := repo.db.Find(&dbUser{Email: login, Password: pass}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (repo *UserRepo) GetById(id uint) (*dbUser, error) {
	var user dbUser

	result := repo.db.Find(&dbUser{ID: id}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
