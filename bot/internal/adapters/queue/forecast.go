package queue

import "bot/internal/domain/entities"

var WeatherCache = make(map[string]map[string]interface{})

func cacheFill(weather entities.Weather) {
	WeatherCache["now"] = map[string]interface{}{
		"datePart":  weather.Forecast.Parts[0].GetPartName(),
		"condition": weather.Fact.GetFactCondition(),
		"temp":      weather.Fact.Temp,
		"feelsLike": weather.Fact.FeelsLike,
	}

	WeatherCache["part1"] = map[string]interface{}{
		"datePart":  weather.Forecast.Parts[0].GetPartName(),
		"condition": weather.Forecast.Parts[0].GetPartCondition(),
		"temp":      weather.Forecast.Parts[0].TempAvg,
		"feelsLike": weather.Forecast.Parts[0].FeelsLike,
	}

	WeatherCache["part2"] = map[string]interface{}{
		"datePart":  weather.Forecast.Parts[1].GetPartName(),
		"condition": weather.Forecast.Parts[1].GetPartCondition(),
		"temp":      weather.Forecast.Parts[1].TempAvg,
		"feelsLike": weather.Forecast.Parts[1].FeelsLike,
	}
}
