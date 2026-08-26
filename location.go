package main

import (
	"fmt"
	"os"
)

const geocodingAPIURL = "https://api.openweathermap.org/geo/1.0/reverse"

type LocationResponse struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

func getLocationName(latitude, longitude float64) (LocationResponse, error) {
	request := APIRequest{
		Latitude:  latitude,
		Longitude: longitude,
		APIKey:    os.Getenv("OPENWEATHER_API_KEY"),
	}

	url := buildAPIURL(geocodingAPIURL, request)

	var locations []LocationResponse

	if err := getJSON(url, &locations); err != nil {
		return LocationResponse{}, fmt.Errorf("get location: %w", err)
	}

	if len(locations) == 0 {
		return LocationResponse{}, fmt.Errorf("location not found")
	}

	return locations[0], nil
}
