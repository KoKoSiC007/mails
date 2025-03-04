package repositories

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type DbRate struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	Name      string
	Rate      float32
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
}

type CurrencyRepo struct {
	db *gorm.DB
}

func NewCurrencyRepo(db *gorm.DB) (*CurrencyRepo, error) {
	return &CurrencyRepo{db: db}, nil
}

func (repo *CurrencyRepo) Get(start, end time.Time) (*[]DbRate, error) {
	var rates []DbRate

	result := repo.db.Where("created_at between @start and @end", sql.Named("start", start), sql.Named("end", end)).Find(&rates)
	if result.Error != nil {
		return nil, result.Error
	}

	return &rates, nil
}
