package main

import (
	"fmt"
)

type APIRequest struct {
	Latitude  float64
	Longitude float64
	APIKey    string
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
