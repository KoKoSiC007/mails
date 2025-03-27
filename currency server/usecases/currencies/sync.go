package currencies

import (
	"time"
)

func (u *CurrencyUseCase) Sync(start, end time.Time, currencies []string) error {
	return u.parser.Parse(start, end, &currencies)
}
