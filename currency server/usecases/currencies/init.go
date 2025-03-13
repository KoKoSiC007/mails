package currencies

import (
	"odyssey/m/v2/repositories"
)

type CurrencyUseCase struct {
	currencyRepo *repositories.CurrencyRepo
	ratesRepo    *repositories.RateRepo
}

func NewCurrencyCase(currencyRepo *repositories.CurrencyRepo, rateRepo *repositories.RateRepo) (*CurrencyUseCase, error) {
	return &CurrencyUseCase{currencyRepo: currencyRepo, ratesRepo: rateRepo}, nil
}
