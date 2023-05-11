package main


import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
	"fmt"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

const tokenlink string = "./config/testtoken.txt"

func main() {
	token, err := readBotToken(tokenlink)
	if err != nil {
		log.Panicf("Token error: ", err)
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	var admin int64 = 193117018

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
	 panic(err)
	}

	passphrase := []byte("myPassphrase")
	
	publicKey := &privateKey.PublicKey
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
				text := "привет!\nэто бот канала «[давай на ты](https://t.me/+eFiRVo8U-NhlNjEy)», и всё, что ты сюда напишешь — анонимно отправится Ксюше 📖\n\nона уже ждёт твои вопросы и истории!"
				msg := tgbotapi.NewMessage(m.Chat.ID, text)
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
				bot.Send(msg)
			} else if m.Chat.ID == admin && m.ReplyToMessage != nil{
				originalmessage := m.ReplyToMessage
				messagetext := originalmessage.Text
				words := strings.Fields(messagetext)

					// Check if there's at least one word
					if len(words) > 0 {
						firstWord := words[0]
						decryptedPlaintext, err := rsa.DecryptOAEP(
							sha256.New(),
							rand.Reader,
							privateKey,
							firstWord,
							passphrase,
						   )
						   
						   if err != nil {
							panic(err)
						   }
						if err != nil {
							// Handle error
							fmt.Println("Error:", err)
							continue
						}
						replychat := string(decryptedPlaintext)
						if m.Text != "" {
							bot.Send(testmsg)
							msg := tgbotapi.NewMessage(replychat, m.Text)
							bot.Send(msg)
						}
					}
			} else {
				passphrase := []byte("myPassphrase")
				chatID, err := rsa.EncryptOAEP(
					sha256.New(),
					rand.Reader,
					publicKey,
					plaintext,
					nil,
				   )
				   
				if err != nil {
					panic(err)
				}
				if m.Text != "" {
					if m.ForwardFrom != nil {
						text := strconv.FormatInt(chatID, 10) + "\n" + m.Text
						msg := tgbotapi.NewMessage(admin, text)
						msg.DisableWebPagePreview = true
						bot.Send(msg)
					} else {
						text := strconv.FormatInt(chatID, 10) + "\n" + m.Text
						msg := tgbotapi.NewMessage(admin, text)
						msg.DisableWebPagePreview = true
						bot.Send(msg)
					}
				}
			}
		}
	}
}
