package main


import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
	"strings"
)

const configlink string = "config.json"
const tokenlink string = "token.txt"

type TempChat struct {
	ChatName      string    `json:"chatName"`
	TimeAdded     int       `json:"timeAdded"` // or you can use time.Duration if you want to store it as a duration instead of a string
	DateTimeAdded time.Time `json:"dateTimeAdded"`
}

type Config struct {
	AllowedUser  int64           `json:"allowedUser"`
	AllowedUsers map[int64]bool  `json:"allowedUsers"`
	Whitelist    map[string]bool `json:"whitelist"`
	AllowedChats map[int64]bool  `json:"allowedChats"`
	TempChats    []TempChat      `json:"tempChats"`
}

type ChatConfig struct{
	ChatID int64
	SuperGroupUsername string
}

func main() {
	config, err := readConfig(configlink)
	if err != nil {
		log.Panicf("Config error: %v", err)
	}

	allowedChats := config.AllowedChats
	whitelist := config.Whitelist

	token, err := readBotToken(tokenlink)
	if err != nil {
		log.Panicf("Token error: ", err)
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	go checkTempChats(bot, &config)

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			m := update.Message
			if m.From == nil {
				continue
			}

			switch update.Message.Command() {
			case "enable_bot":
				handleEnableBotCommand(bot, m, &config)
			case "disable_bot":
				handleDisableBotCommand(bot, m, &config)
			case "whitelist":
				handleWhitelistCommand(bot, m, &config)
			case "tempadd":
				handleTempAddCommand(bot, m, &config)
			case "list":
				handleListCommand(bot, m, &config)
			case "templist":
				handleListTempChatsCommand(bot, m, &config)
			case "remove":
				handleRemoveCommand(bot, m, &config)
			}

			err := handleChannelLink(bot, m, &config)

			if err != nil {
				log.Printf("Error deleting message: ", err)
			}

			if m.SenderChat != nil && m.SenderChat.Type == "channel" && allowedChats[m.Chat.ID] && !whitelist[strings.ToLower(m.SenderChat.UserName)] {
				_, err := bot.Send(tgbotapi.DeleteMessageConfig{
					ChatID:    update.Message.Chat.ID,
					MessageID: update.Message.MessageID,
				})
				if err != nil {
					log.Printf("[%s] %s,   err: %s", update.Message.From.UserName, update.Message.Text, err.Error())
					continue
				}
			}
		}
	}
}
