package queue

import (
	"io/ioutil"
	"net/http"

	"scheduler/internal/config"
)

func (a *App) YandexForecastCall(cfg *config.Config) ([]byte, error) {
	url := "https://api.weather.yandex.ru/v2/informers?lat=55.75222&lon=37.61556"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("X-Yandex-API-Key", cfg.YandexApiToken)
	resp, err := http.DefaultClient.Do(request)
	defer func() { _ = resp.Body.Close() }()

	body, _ := ioutil.ReadAll(resp.Body)

	//log.Println("call ya")

	return body, nil
}
