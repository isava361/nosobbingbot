package main


import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
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
				text := "привет!\nэто бот канала «[давай на ты](https://t.me/+eFiRVo8U-NhlNjEy)», и всё, что ты сюда напишешь — анонимно отправится Ксюше 📖\n\nона уже ждёт твои вопросы и истории!"
				msg := tgbotapi.NewMessage(m.Chat.ID, text)
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true
				bot.Send(msg)
				continue
			} else if m.Chat.Type == "private" && m.Chat.ID != admin {
				
				chatID := strconv.FormatInt(encryptID(m.Chat.ID),10)
				if m.Text != "" {
					text := chatID + "\n" + m.Text
					msg := tgbotapi.NewMessage(admin, text)
					bot.Send(msg)
				}else if m.Photo != nil {
					text := chatID +":\n" + m.Caption
					existingFileID := m.Photo[len(m.Photo)-1].FileID
					photoMsg := tgbotapi.NewPhoto(admin, tgbotapi.FileID(existingFileID))
					photoMsg.Caption = text
					bot.Send(photoMsg)
				} else if m.Audio != nil {
					text := chatID +":\n" + m.Caption
					existingFileID := m.Audio.FileID
					audioMsg := tgbotapi.NewAudio(admin, tgbotapi.FileID(existingFileID))
					audioMsg.Caption = text
					bot.Send(audioMsg)
				} else if m.Document != nil {
					text := chatID +":\n" + m.Caption
					existingFileID := m.Document.FileID
					documentMsg := tgbotapi.NewDocument(admin, tgbotapi.FileID(existingFileID))
					documentMsg.Caption = text
					bot.Send(documentMsg)
				} else if m.Video != nil {
					text := chatID +":\n" + m.Caption
					existingFileID := m.Video.FileID
					videoMsg := tgbotapi.NewVideo(admin, tgbotapi.FileID(existingFileID))
					videoMsg.Caption = text
					bot.Send(videoMsg)
				} else if m.Voice != nil {
					text := chatID + ":"
					msg := tgbotapi.NewMessage(admin, text)
					existingFileID := m.Voice.FileID
					voiceMsg := tgbotapi.NewVoice(admin, tgbotapi.FileID(existingFileID))
					bot.Send(msg)
					bot.Send(voiceMsg)
				} else if m.VideoNote != nil {
					text := chatID + ":"
					msg := tgbotapi.NewMessage(admin, text)
					existingFileID := m.VideoNote.FileID
					length := m.VideoNote.Length
					videoNoteMsg := tgbotapi.NewVideoNote(admin, length, tgbotapi.FileID(existingFileID))
					bot.Send(msg)
					bot.Send(videoNoteMsg)
				}
			} else if m.Chat.ID == admin && m.ReplyToMessage != nil {
				originalmessage := m.ReplyToMessage
				messagetext := originalmessage.Text
				if m.Caption != "" {
					messagetext = originalmessage.Text
				}
				words := strings.Fields(messagetext)
				if len(words) > 0 {
					firstWord := words[0]
					encryptedmessage, _ := strconv.ParseInt(string(firstWord), 10, 64)
					replychat := decryptID(encryptedmessage)

					if m.Text != "" {
						msg := tgbotapi.NewMessage(replychat, m.Text)
						bot.Send(msg)
					} else if m.Photo != nil {
						existingFileID := m.Photo[len(m.Photo)-1].FileID
						photoMsg := tgbotapi.NewPhoto(replychat, tgbotapi.FileID(existingFileID))
						photoMsg.Caption = m.Caption
						bot.Send(photoMsg)
					} else if m.Audio != nil {
						existingFileID := m.Audio.FileID
						audioMsg := tgbotapi.NewAudio(replychat, tgbotapi.FileID(existingFileID))
						audioMsg.Caption = m.Caption
						bot.Send(audioMsg)
					} else if m.Document != nil {
						existingFileID := m.Document.FileID
						documentMsg := tgbotapi.NewDocument(replychat, tgbotapi.FileID(existingFileID))
						documentMsg.Caption = m.Caption
						bot.Send(documentMsg)
					} else if m.Video != nil {
						existingFileID := m.Video.FileID
						videoMsg := tgbotapi.NewVideo(replychat, tgbotapi.FileID(existingFileID))
						videoMsg.Caption = m.Caption
						bot.Send(videoMsg)
					} else if m.Voice != nil {
						existingFileID := m.Voice.FileID
						voiceMsg := tgbotapi.NewVoice(replychat, tgbotapi.FileID(existingFileID))
						bot.Send(voiceMsg)
					} else if m.VideoNote != nil {
						existingFileID := m.VideoNote.FileID
						length := m.VideoNote.Length
						videoNoteMsg := tgbotapi.NewVideoNote(replychat, length, tgbotapi.FileID(existingFileID))
						bot.Send(videoNoteMsg)
					}
				}
			} 
		}
	}
}
