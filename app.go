package main


import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
//	"strings"
)

const tokenlink string = "./config/token.txt"

func main() {
	token, err := readBotToken(tokenlink)
	if err != nil {
		log.Panicf("Token error: ", err)
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	var admin int64 = 852084868
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			m := update.Message
			if m.From == nil {
				continue
			}
			if m.Text == "/start" {
				text := "привет!/n это бот канала «давай на ты», и всё, что ты сюда напишешь — анонимно отправится Ксюше 📖 /n она уже ждёт твои вопросы и истории!"
				msg := tgbotapi.NewMessage(m.Chat.ID, text)
				bot.Send(msg)
			} else if m.Chat.ID == admin {
				continue
			} else {
				text := m.Text
				msg := tgbotapi.NewMessage(admin, text)
				bot.Send(msg)
			}
		}
	}
}
