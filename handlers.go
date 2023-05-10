package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
	"time"
	"fmt"
)

func handleListCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, config *Config) {
	allowedUsers := config.AllowedUsers
	if allowedUsers[message.From.ID] {
		var whitelistLines []string
		for channel := range config.Whitelist {
			whitelistLines = append(whitelistLines, channel)
		}
		whitelistString := strings.Join(whitelistLines, "\n")
		if whitelistString == "" {
			whitelistString = "The whitelist is empty."
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, "Whitelist contents: \n\n"+whitelistString)
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "You are not authorized to use this command.")
		bot.Send(msg)
	}
}

func handleListTempChatsCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, config *Config) {
	allowedUsers := config.AllowedUsers
	tempChats := config.TempChats

	if allowedUsers[message.From.ID] {
		if len(tempChats) == 0 {
			msg := tgbotapi.NewMessage(message.Chat.ID, "No temporary chats found.")
			bot.Send(msg)
			return
		}

		tempChatList := ""

		for _, tempChat := range tempChats {
			tempChatList += fmt.Sprintf("%s - %dm - %v\n", tempChat.ChatName, tempChat.TimeAdded, tempChat.DateTimeAdded.Format("2006-01-02 15:04:05"))
		}

		msg := tgbotapi.NewMessage(message.Chat.ID, "Temporary added chats:\n\n"+tempChatList)
		bot.Send(msg)

	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "You are not authorized to use this command.")
		bot.Send(msg)
	}
}

func handleTempAddCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, config *Config) {
	if config.AllowedUsers[message.From.ID] {
		args := strings.Split(strings.ToLower(message.Text), " ")
		if len(args) < 3 {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Usage: /tempadd @channelusername duration_in_minutes")
			_, _ = bot.Send(msg)
			return
		}

		duration, err := strconv.Atoi(args[2])
		if err != nil || duration <= 0 || duration > 30 {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Duration must be an integer between 1 and 30 minutes")
			_, _ = bot.Send(msg)
			return
		}

		channel := args[1]

		if channel != "" {
			config.Whitelist[channel] = true
			err := addChannelToWhitelist(config, channel)
			if err != nil {
				log.Printf("Error updating whitelist: %v", err)
			} else {
				msg := tgbotapi.NewMessage(message.Chat.ID, "Channel temporarily added to whitelist: "+channel)
				bot.Send(msg)

				// Add temporary chat to config and update JSON
				addTempChat(config, channel, duration, time.Now())
			}
		}
	} else if !config.AllowedUsers[message.From.ID] {
		msg := tgbotapi.NewMessage(message.Chat.ID, "You are not authorized to use this command")
		_, _ = bot.Send(msg)
		return
	}
}



func handleRemoveCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, config *Config) {
	if message.From.ID == config.AllowedUser {
		channel := strings.TrimSpace(strings.TrimPrefix(strings.ToLower(message.Text), "/remove"))
		if channel != "" {
			err := removeFromWhitelist(config, channel)
			if err != nil {
				log.Printf("Error removing channel from whitelist: %v", err)
			} else {
				removeFromWhitelist(config, channel)
				msg := tgbotapi.NewMessage(message.Chat.ID, "Channel removed from whitelist: "+channel)
				bot.Send(msg)
			}
		}
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "You are not authorized to use this command")
		bot.Send(msg)
	}
}

func handleWhitelistCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, config *Config) {
	allowedUser := config.AllowedUser
	whitelist := config.Whitelist
	if message.From.ID == allowedUser {
		channel := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(message.Text, "/whitelist")))
		if channel != "" {
			whitelist[channel] = true
			config.Whitelist = whitelist
			err := writeConfig(configlink, *config)
			if err != nil {
				log.Printf("Error updating whitelist: %v", err)
			} else {
				msg := tgbotapi.NewMessage(message.Chat.ID, "Channel added to whitelist: "+channel)
				bot.Send(msg)
			}
		}
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "You are not authorized to use this command")
		bot.Send(msg)
	}
}

func handleEnableBotCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, config *Config) {
	allowedUsers := config.AllowedUsers
	if message.Chat.Type == "supergroup" && allowedUsers[message.From.ID] {
		err := addChatToAllowedChats(config, message.Chat.ID)
		if err != nil {
			log.Printf("Error updating allowed chats: %v", err)
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Bot enabled in this chat.")
			bot.Send(msg)
		}
	}
}

func handleDisableBotCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message, config *Config) {
	allowedUsers := config.AllowedUsers
	allowedChats := config.AllowedChats
	if message.Chat.Type == "supergroup" && allowedUsers[message.From.ID] {
		allowedChats[message.Chat.ID] = true
		err := removeChatFromAllowedChats(config, message.Chat.ID)
		if err != nil {
			log.Printf("Error updating allowed chats: %v", err)
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Bot disabled in this chat.")
			bot.Send(msg)
		}
	}
}


func handleChannelLink(bot *tgbotapi.BotAPI, message *tgbotapi.Message, config *Config) error {
 whitelist := config.Whitelist

	containsChannelLink := getChannelUsername(bot, message.Text)
	if (message.SenderChat != nil && !whitelist[strings.ToLower(message.SenderChat.UserName)]) || (message.SenderChat == nil &&!whitelist[strings.ToLower(message.From.UserName)]){
		if containsChannelLink {
				fmt.Println("Tryingtodeletemessage")
				_, err := bot.Send(tgbotapi.DeleteMessageConfig{
					ChatID:    message.Chat.ID,
					MessageID: message.MessageID,
				})
				if err != nil {
					log.Printf("[%s] %s,   err: %s", message.From.UserName, message.Text, err.Error())
					return err
				}
		}
	}
	
	return nil
}


func getChannelUsername(bot *tgbotapi.BotAPI, messageText string) (bool) {
	fmt.Println("Is there a link here?")
	fmt.Println(messageText)

	args := strings.Split(strings.ToLower(messageText), " ")

	for _, arg := range args {
		if strings.HasPrefix(arg, "@") {
			channelUsername := arg[:]
			chat, _ := bot.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{SuperGroupUsername: channelUsername}})
			if chat.Type == "channel"{
				return true
			}
		} else if strings.HasPrefix(arg, "t.me/") {
			channelUsername := "@" + arg[len("tm.me/"):]
			chat, _ := bot.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{SuperGroupUsername: channelUsername}})
			if chat.Type == "channel"{
				return true
			}
		} else if strings.HasPrefix(arg, "telegram.me/") {
			channelUsername := "@" + arg[len("telegram.me/"):]
			chat, _ := bot.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{SuperGroupUsername: channelUsername}})
			if chat.Type == "channel"{
				return true
			}
		}
	}

	return false
}



