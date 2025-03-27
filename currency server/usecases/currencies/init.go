package currencies

import (
	"odyssey/m/v2/internal/parser"
	"odyssey/m/v2/repositories"
)

type CurrencyUseCase struct {
	currencyRepo *repositories.CurrencyRepo
	ratesRepo    *repositories.RateRepo
	parser       *parser.Parser
}

func NewCurrencyCase(currencyRepo *repositories.CurrencyRepo, rateRepo *repositories.RateRepo) (*CurrencyUseCase, error) {
	parser := parser.NewParser(rateRepo)

	return &CurrencyUseCase{currencyRepo: currencyRepo, ratesRepo: rateRepo, parser: parser}, nil
}
