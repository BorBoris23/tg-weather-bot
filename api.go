package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIRequest struct {
	Latitude  float64
	Longitude float64
	APIKey    string
}

func getJSON(url string, result interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

func buildAPIURL(baseURL string, request APIRequest) string {
	return fmt.Sprintf(
		"%s?lat=%f&lon=%f&appid=%s",
		baseURL,
		request.Latitude,
		request.Longitude,
		request.APIKey,
	)
}
