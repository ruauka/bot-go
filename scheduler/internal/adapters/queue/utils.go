package queue

import (
	"bytes"
	"log"
	"time"
)

var ReqBodyBytes = new(bytes.Buffer)

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

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
