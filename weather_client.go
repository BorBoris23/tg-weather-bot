package main

import (
	"fmt"
	"net/http"
)

const geocodingAPIURL = "https://api.openweathermap.org/geo/1.0/reverse"
const weatherAPIURL = "https://api.openweathermap.org/data/2.5/weather"

type LocationResponse struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

type WeatherClient struct {
	httpClient *http.Client
	apiKey     string
}

func NewWeatherClient(httpClient *http.Client, apiKey string) *WeatherClient {
	return &WeatherClient{
		httpClient: httpClient,
		apiKey:     apiKey,
	}
}

func (w *WeatherClient) getWeather(latitude, longitude float64) (WeatherResponse, error) {
	request := APIRequest{
		Latitude:  latitude,
		Longitude: longitude,
		APIKey:    w.apiKey,
	}

	url := buildAPIURL(weatherAPIURL, request)

	var weather WeatherResponse

	if err := getJSON(w.httpClient, url, &weather); err != nil {
		return WeatherResponse{}, fmt.Errorf("get weather: %w", err)
	}

	return weather, nil
}

func (w *WeatherClient) getLocationName(latitude, longitude float64) (LocationResponse, error) {
	request := APIRequest{
		Latitude:  latitude,
		Longitude: longitude,
		APIKey:    w.apiKey,
	}

	url := buildAPIURL(geocodingAPIURL, request)

	var locations []LocationResponse

	if err := getJSON(w.httpClient, url, &locations); err != nil {
		return LocationResponse{}, fmt.Errorf("get location: %w", err)
	}

	if len(locations) == 0 {
		return LocationResponse{}, fmt.Errorf("location not found")
	}

	return locations[0], nil
}
