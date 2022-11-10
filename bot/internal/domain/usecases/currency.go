package usecases

import (
	"fmt"
	"time"

	"github.com/matperez/go-cbr-client"
)

type currencyUsecase struct {
	client cbr.Client
}

func NewCurrencyUsecase(CBRFClient cbr.Client) CurrencyUsecase {
	return &currencyUsecase{
		client: CBRFClient,
	}
}

func (c *currencyUsecase) Get(currency string, date time.Time) (float64, error) {
	rate, err := c.client.GetRate(currency, date)
	if err != nil {
		return 0, fmt.Errorf("connect cb error")
	}

	return rate, nil
}
