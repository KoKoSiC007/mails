package currencies

import (
	"odyssey/m/v2/models"
	"odyssey/m/v2/repositories"
	"time"
)

func (u *CurrencyUseCase) GetRates(start, end time.Time, currencies []string) (*[]models.CurrencyInner, error) {
	result, err := u.ratesRepo.Get(start, end, currencies)
	if err != nil {
		return nil, err
	}

	curMap := u.mapRates(result)
	rates := u.getReport(&curMap)

	return rates, nil
}

func (u *CurrencyUseCase) mapRates(rates *[]repositories.DbRate) map[string][]repositories.DbRate {
	var curMap = make(map[string][]repositories.DbRate)

	for _, val := range *rates {
		curMap[val.Name] = append(curMap[val.Name], val)
	}

	return curMap
}

func (u *CurrencyUseCase) getReport(curMap *map[string][]repositories.DbRate) *[]models.CurrencyInner {
	var currencies []models.CurrencyInner
	for name, rates := range *curMap {
		var reportPart = models.CurrencyInner{Name: name}
		var min, max, avg, summ float32
		min = rates[0].Rate

		for _, rate := range rates {
			if min > rate.Rate {
				min = rate.Rate
			}
			if max < rate.Rate {
				max = rate.Rate
			}
			summ += rate.Rate
		}
		avg = summ / float32(len(rates))

		reportPart.AvgRate = avg
		reportPart.MaxRate = max
		reportPart.MinRate = min

		currencies = append(currencies, reportPart)
	}

	return &currencies
}
