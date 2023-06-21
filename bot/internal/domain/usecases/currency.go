package usecases

import (
	"fmt"
	"time"

	"github.com/ivangurin/cbrf-go"
)

type currencyUsecase struct{}

func NewCurrencyUsecase() CurrencyUsecase { return &currencyUsecase{} }

func (c *currencyUsecase) Get(currency string, date time.Time) (float64, error) {
	rate, err := cbrf.GetExchangeRate(currency, date)
	if err != nil {
		return 0, fmt.Errorf("connect cb error")
	}

	return rate, nil
}
