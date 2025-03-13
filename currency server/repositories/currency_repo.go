package repositories

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type dbCurrency struct {
	gorm.Model
	ID        uint      `gorm:"privateKey"`
	Name      string    `gorm:"not null"`
	Schedule  string    `gorm:"not null"`
	Enable    bool      `gorm:"not null, default:true"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type CurrencyRepo struct {
	db *gorm.DB
}

func NewCurrencyRepo(db *gorm.DB) (*CurrencyRepo, error) {
	err := db.AutoMigrate(&dbCurrency{})
	if err != nil {
		fmt.Println(err.Error())
	}

	return &CurrencyRepo{db: db}, nil
}

func (repo *CurrencyRepo) Get() (*[]dbCurrency, error) {
	var currencies []dbCurrency

	result := repo.db.Find(&currencies)
	if result.Error != nil {
		return nil, result.Error
	}

	return &currencies, nil
}
