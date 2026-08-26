package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Location == nil {
			continue
		}

		latitude := update.Message.Location.Latitude
		longitude := update.Message.Location.Longitude

		weather, err := getWeather(latitude, longitude)
		if err != nil {
			log.Println(err)
			continue
		}

		location, err := getLocationName(latitude, longitude)
		if err != nil {
			log.Printf("get location name: %v", err)
			continue
		}

		msgText := formatWeather(weather, location)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)

		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}
