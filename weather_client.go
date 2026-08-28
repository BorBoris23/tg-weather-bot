package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const openWeatherBaseURL = "https://api.openweathermap.org"

const (
	weatherEndpoint   = "/data/2.5/weather"
	geocodingEndpoint = "/geo/1.0/reverse"
)

type APIRequest struct {
	Latitude  float64
	Longitude float64
	APIKey    string
}

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

	url := buildAPIURL(weatherEndpoint, request)

	var weather WeatherResponse

	if err := w.getJSON(url, &weather); err != nil {
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

	url := buildAPIURL(geocodingEndpoint, request)

	var locations []LocationResponse

	if err := w.getJSON(url, &locations); err != nil {
		return LocationResponse{}, fmt.Errorf("get location: %w", err)
	}

	if len(locations) == 0 {
		return LocationResponse{}, fmt.Errorf("location not found")
	}

	return locations[0], nil
}

func (w *WeatherClient) getJSON(url string, result interface{}) error {
	response, err := w.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("request API: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

func buildAPIURL(endpoint string, request APIRequest) string {
	return fmt.Sprintf(
		"%s%s?lat=%f&lon=%f&appid=%s",
		openWeatherBaseURL,
		endpoint,
		request.Latitude,
		request.Longitude,
		request.APIKey,
	)
}
