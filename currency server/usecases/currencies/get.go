package currencies

import "odyssey/m/v2/models/currency"

func (u *CurrencyUseCase) Get() (*[]currency.Currency, error) {
	result, err := u.currencyRepo.Get()
	if err != nil {
		return nil, err
	}
	var currencies []currency.Currency
	for _, val := range *result {
		currency := currency.Currency{
			Id:       val.ID,
			Name:     val.Name,
			Enable:   val.Enable,
			Schedule: val.Schedule,
		}

		currencies = append(currencies, currency)
	}

	return &currencies, nil
}
