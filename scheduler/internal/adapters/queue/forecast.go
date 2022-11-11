package queue

import (
	"io/ioutil"
	"log"
	"net/http"
)

func (a *App) YandexForecastCall() []byte {
	url := "https://api.weather.yandex.ru/v2/informers?lat=55.75222&lon=37.61556"

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to connect to Yandex: %s", err.Error())
	}

	request.Header.Add("X-Yandex-API-Key", "03633f27-bab2-482e-b8f4-532f7775c13c")
	resp, err := http.DefaultClient.Do(request)
	defer func() { _ = resp.Body.Close() }()

	body, _ := ioutil.ReadAll(resp.Body)

	//log.Println("call ya")

	return body
}
