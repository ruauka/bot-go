package queue

import (
	"bytes"
	"time"
)

var (
	ReqBodyBytes = new(bytes.Buffer)
	deduct       = 3
)

// Layout - шаблон формата даты для изменения типа string -> time.Time.
const Layout = "02.01.2006 15:04"

// convertDate - преобразование входящий даты-строки в дату-time.Time.
func convertDate(date string) time.Time {
	parseDate, err := time.Parse(Layout, date)
	if err != nil {
		return time.Time{}
	}

	return parseDate
}

func morningCheck(date string) bool {
	currentDate := convertDate(date)
	if currentDate.Year() == time.Now().Year() &&
		currentDate.Month() == time.Now().Month() &&
		currentDate.Day() == time.Now().Day() &&
		time.Now().Hour() == 7-deduct && time.Now().Minute() == 40 &&
		(time.Now().Second() == 0 ||
			time.Now().Second() == 1 ||
			time.Now().Second() == 2 ||
			time.Now().Second() == 3 ||
			time.Now().Second() == 4) {
		return true
	}

	return false
}

func eveningCheck(date string) bool {
	currentDate := convertDate(date)
	if currentDate.Year() == time.Now().Year() &&
		currentDate.Month() == time.Now().Month() &&
		currentDate.Day() == time.Now().Day()+1 &&
		time.Now().Hour() == 21-deduct && time.Now().Minute() == 30 &&
		(time.Now().Second() == 0 ||
			time.Now().Second() == 1 ||
			time.Now().Second() == 2 ||
			time.Now().Second() == 3 ||
			time.Now().Second() == 4) {
		return true
	}

	return false
}
