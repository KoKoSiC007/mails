package repositories

import (
	"time"

	"example.com/m/v2/internal/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dbRate struct {
	gorm.Model
	ID        uint `gorm:"primaryKey"`
	Name      string
	Rate      float32
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Repo struct {
	db *gorm.DB
}

func NewRepo() (*Repo, error) {
	dsn := "host=localhost user=postgres password=234492 dbname=currencies port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&dbRate{})

	return &Repo{db: db}, nil
}

func (repo *Repo) Create(rate entities.Rate) (*dbRate, error) {
	entity := dbRate{Name: rate.Name, Rate: rate.Rate}
	result := repo.db.Create(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	return &entity, nil
}
