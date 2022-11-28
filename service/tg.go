package service

import (
	"encoding/json"
	"io"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/neversi/baiterek-app-api/config"
)

type TGBot struct {
	io.Closer
	config.TGBotConfig
	bot         *tgbotapi.BotAPI
	authStorage *Storage[Authorization]
	Logins      chan Login
	done        chan struct{}
	end         chan struct{}
}

func NewBot(cfg *config.TGBotConfig) *TGBot {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 0

	updates := bot.GetUpdatesChan(u)
	done := make(chan struct{}, 1)
	end := make(chan struct{}, 1)
	tb := &TGBot{
		TGBotConfig: *cfg,
		Logins:      make(chan Login, 100),
		authStorage: NewStorage[Authorization](),
		bot:         bot,
		done:        done,
		end:         end,
	}
	go func() {
		defer func() {
			end <- struct{}{}
		}()
		for {
			select {
			case <-done:
				log.Println("closing reading chan")
				return
			case update := <-updates:
				if tb.authorize(update) {
					log.Println("not authorized skip...")
				}
			}
		}
	}()

	go tb.Listen()

	return tb
}

func (tb *TGBot) Close() error {
	tb.done <- struct{}{}
	<-tb.end

	tb.done <- struct{}{}
	<-tb.end

	log.Println("ending bot")
	return nil
}

func (tb *TGBot) IsAuthorized(username string) bool {
	return tb.authStorage.Get(username) != nil
}

const (
	needAuth = "You need to authorize, please pass the unique key: "
	wrongKey = "Wrong key"
)

func (tb *TGBot) authorize(update tgbotapi.Update) bool {
	if val := tb.authStorage.Get(update.Message.From.UserName); val != nil && val.Authorized {
		return false
	}

	if update.Message.Text != "/start" {
		if update.Message.Text != tb.Key {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, wrongKey)
			msg.ReplyToMessageID = update.Message.MessageID

			tb.bot.Send(msg)
			return true
		}
		tb.authStorage.Set(update.Message.From.UserName, Authorization{
			ChatID:     update.Message.Chat.ID,
			Authorized: true,
		})
		return false
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, needAuth)
	msg.ReplyToMessageID = update.Message.MessageID

	tb.bot.Send(msg)

	tb.authStorage.Set(update.Message.From.UserName, Authorization{
		ChatID:     update.Message.Chat.ID,
		Authorized: false,
	})

	return true
}

func (tb *TGBot) Listen() {
	go func() {
		defer func() {
			tb.end <- struct{}{}
		}()
		for {
			select {
			case <-tb.done:
				log.Println("ending listener of logins")
				close(tb.Logins)
				return
			case login := <-tb.Logins:
				log.Printf("login data: %v\n", login)
				body, err := json.Marshal(login)
				if err != nil {
					log.Printf("error: %v\n", err)
					continue
				}
				for _, val := range tb.authStorage.List() {
					if val.Authorized {
						msg := tgbotapi.NewMessage(val.ChatID, string(body))
						tb.bot.Send(msg)
					}
				}
			}
		}
	}()
}
