package main

import (
	"log"
	"math/rand"
	"sync/atomic"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
)

const (
	greetingMsg = "Привет! Я Юрий Ткаченко. Теперь я буду следить за нравственным порядком и ущемлением прав LGBTQIA здесь и по-хорошему просить удалить нахуй этот позор."
	badJokeMsg  = "Мне кажется, что вы не смешно пошутили.\nПредлагаю вам удалить нахуй этот позор."
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	var messagesCount uint64

	// config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Panic(err)
	}

	token := viper.GetString("token")
	entropy := viper.GetInt("entropy")
	repeatEvery := viper.GetInt("repeat")

	// bot setup
	bot, err := tg.NewBotAPI(token)

	if err != nil {
		log.Panic(err)
	}

	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		switch {
		case update.Message != nil:
			go func(msg *tg.Message) {
				if msg.IsCommand() {
					switch msg.Command() {
					case "start":
						reply := tg.NewMessage(msg.Chat.ID, greetingMsg)
						bot.Send(reply)
						return
					}
				}

				if len(msg.Text) == 0 {
					return
				}

				count := atomic.AddUint64(&messagesCount, 1)

				if count >= uint64(repeatEvery) && count >= uint64(repeatEvery)+uint64(rand.Intn(entropy)) {
					atomic.StoreUint64(&messagesCount, 0)
					reply := tg.NewMessage(msg.Chat.ID, badJokeMsg)
					reply.ReplyToMessageID = msg.MessageID

					if _, err := bot.Send(reply); err != nil {
						log.Println(err)
					}
				}
			}(update.Message)
		}
	}
}
