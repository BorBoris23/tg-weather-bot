package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	weatherClient := NewWeatherClient(
		httpClient,
		os.Getenv("OPENWEATHER_API_KEY"),
	)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if err := handleUpdate(bot, weatherClient, update); err != nil {
			log.Printf("handle update: %v", err)
		}
	}
}

func handleUpdate(bot *tgbotapi.BotAPI, weatherClient *WeatherClient, update tgbotapi.Update) error {
	if update.Message == nil || update.Message.Location == nil {
		return nil
	}

	latitude := update.Message.Location.Latitude
	longitude := update.Message.Location.Longitude

	weather, err := weatherClient.getWeather(latitude, longitude)
	if err != nil {
		return fmt.Errorf("get weather: %w", err)
	}

	location, err := weatherClient.getLocationName(latitude, longitude)
	if err != nil {
		return fmt.Errorf("get location name: %w", err)
	}

	msgText := formatWeather(weather, location)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)

	if _, err := bot.Send(msg); err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
