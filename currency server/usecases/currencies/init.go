package currencies

import (
	"odyssey/m/v2/repositories"
)

type CurrencyUseCase struct {
	currencyRepo *repositories.CurrencyRepo
}

func NewCurrencyCase(repo *repositories.CurrencyRepo) (*CurrencyUseCase, error) {
	return &CurrencyUseCase{currencyRepo: repo}, nil
}
