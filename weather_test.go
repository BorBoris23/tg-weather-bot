package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBuildAPIURL(t *testing.T) {
	request := APIRequest{
		Latitude:  68.882393,
		Longitude: 81.000045,
		APIKey:    "test-key",
	}

	got := buildAPIURL("https://example.com/weather", request)

	expected := "https://example.com/weather?lat=68.882393&lon=81.000045&appid=test-key"

	if got != expected {
		t.Errorf("unexpected URL:\nwant: %s\ngot:  %s", expected, got)
	}
}

func TestFormatWeather(t *testing.T) {
	weather := WeatherResponse{
		Name: "Amsterdam",
		Metrics: Metrics{
			Temp:      23.39,
			FeelsLike: 22.75,
			Humidity:  37,
		},
		Weather: []Weather{
			{
				Description: "overcast clouds",
			},
		},
		Wind: Wind{
			Speed: 8.29,
		},
	}

	location := LocationResponse{
		Name:    "Amsterdam",
		Country: "NL",
	}

	got := formatWeather(weather, location)

	expected := "Weather in Amsterdam, NL:\n\nTemperature: 23.4°C\nFeels like: 22.8°C\nCondition: overcast clouds\nHumidity: 37%\nWind speed: 8.3 m/s"

	if got != expected {
		t.Errorf("unexpected result:\nwant: %s\ngot:  %s", expected, got)
	}
}

func TestGetLocationNameNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, `[]`)
	}))
	defer server.Close()

	var locations []LocationResponse

	err := getJSON(server.URL, &locations)
	if err != nil {
		t.Fatalf("getJSON returned error: %v", err)
	}

	if len(locations) != 0 {
		t.Fatalf("expected no locations, got %d", len(locations))
	}
}

func TestGetLocationName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, `[
			{
				"name": "Amsterdam",
				"country": "NL"
			}
		]`)
	}))
	defer server.Close()

	var locations []LocationResponse

	err := getJSON(server.URL, &locations)
	if err != nil {
		t.Fatalf("getJSON returned error: %v", err)
	}

	if len(locations) != 1 {
		t.Fatalf("expected 1 location, got %d", len(locations))
	}

	if locations[0].Name != "Amsterdam" {
		t.Errorf("unexpected location name: %s", locations[0].Name)
	}

	if locations[0].Country != "NL" {
		t.Errorf("unexpected country: %s", locations[0].Country)
	}
}

func TestGetJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		fmt.Fprint(w, `{
			"name": "Amsterdam",
			"main": {
				"temp": 23.39,
				"feels_like": 22.75,
				"humidity": 37
			},
			"weather": [
				{
					"description": "overcast clouds"
				}
			],
			"wind": {
				"speed": 8.29
			}
		}`)
	}))
	defer server.Close()

	var weather WeatherResponse

	err := getJSON(server.URL, &weather)
	if err != nil {
		t.Fatalf("getJSON returned error: %v", err)
	}

	if weather.Name != "Amsterdam" {
		t.Errorf("unexpected name: %s", weather.Name)
	}

	if weather.Metrics.Temp != 23.39 {
		t.Errorf("unexpected temperature: %f", weather.Metrics.Temp)
	}

	if weather.Metrics.FeelsLike != 22.75 {
		t.Errorf("unexpected feels like: %f", weather.Metrics.FeelsLike)
	}

	if weather.Metrics.Humidity != 37 {
		t.Errorf("unexpected humidity: %d", weather.Metrics.Humidity)
	}

	if len(weather.Weather) == 0 {
		t.Fatal("weather array is empty")
	}

	if weather.Weather[0].Description != "overcast clouds" {
		t.Errorf("unexpected description: %s", weather.Weather[0].Description)
	}

	if weather.Wind.Speed != 8.29 {
		t.Errorf("unexpected wind speed: %f", weather.Wind.Speed)
	}
}

func TestGetJSONInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `invalid json`)
	}))
	defer server.Close()

	var weather WeatherResponse

	err := getJSON(server.URL, &weather)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "decode response") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestGetJSONStatusError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}))
	defer server.Close()

	var weather WeatherResponse

	err := getJSON(server.URL, &weather)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "API returned status 401") {
		t.Errorf("unexpected error: %v", err)
	}
}
