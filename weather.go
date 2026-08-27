package main

import (
	"fmt"
)

type Metrics struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	Humidity  int     `json:"humidity"`
}

type Weather struct {
	Description string `json:"description"`
}

type Wind struct {
	Speed float64 `json:"speed"`
}

type WeatherResponse struct {
	Name    string    `json:"name"`
	Metrics Metrics   `json:"main"`
	Weather []Weather `json:"weather"`
	Wind    Wind      `json:"wind"`
}

func formatWeather(weather WeatherResponse, location LocationResponse) string {
	return fmt.Sprintf(
		"Weather in %s, %s:\n\nTemperature: %.1f°C\nFeels like: %.1f°C\nCondition: %s\nHumidity: %d%%\nWind speed: %.1f m/s",
		location.Name,
		location.Country,
		weather.Metrics.Temp,
		weather.Metrics.FeelsLike,
		weather.Weather[0].Description,
		weather.Metrics.Humidity,
		weather.Wind.Speed,
	)
}
