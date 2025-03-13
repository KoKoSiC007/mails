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

type RateRepo struct {
	db *gorm.DB
}

func NewRateRepo(db *gorm.DB) (*RateRepo, error) {
	return &RateRepo{db: db}, nil
}

func (repo *RateRepo) Get(start, end time.Time) (*[]DbRate, error) {
	var rates []DbRate

	result := repo.db.Where("created_at between @start and @end", sql.Named("start", start), sql.Named("end", end)).Find(&rates)
	if result.Error != nil {
		return nil, result.Error
	}

	return &rates, nil
}
