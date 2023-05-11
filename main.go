package main

import (
	"log"

	"github.com/expose443/BookLinkerBot/client"
	"github.com/expose443/BookLinkerBot/token"
)

const BotApi = "https://api.telegram.org/bot"

func main() {
	botToken, err := token.FromEnv("BOT_TOKEN")
	if err != nil {
		log.Println(err)
		return
	}
	botUrl := BotApi + botToken
	client := client.NewHttpClient(botUrl)
	botName, err := token.FromEnv("BOT_NAME")
	if err != nil {
		botName = "Bot_Name"
	}
	log.Println("Bot running... name in telegram:", botName)
	for {
		updates, err := client.GetUpdates()
		if err != nil {
			log.Println(err)

			return
		}
		for _, update := range updates {

			err = client.Respond(update)
			if err != nil {
				log.Println(err)
				return
			}
			client.Offset = update.UpdateId + 1
		}

	}
}
