package main

import (
	"fmt"
	"os"
)

const weatherAPIURL = "https://api.openweathermap.org/data/2.5/weather"

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

func getWeather(latitude, longitude float64) (WeatherResponse, error) {
	request := APIRequest{
		Latitude:  latitude,
		Longitude: longitude,
		APIKey:    os.Getenv("OPENWEATHER_API_KEY"),
	}

	url := buildAPIURL(weatherAPIURL, request)

	var weather WeatherResponse

	if err := getJSON(url, &weather); err != nil {
		return WeatherResponse{}, fmt.Errorf("get weather: %w", err)
	}

	return weather, nil
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
