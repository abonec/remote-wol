package main

import (
	"os"
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

var (
	bot *tgbotapi.BotAPI
)

func startTelegramBot() {
	var tb_key string
	if tb_key = os.Getenv("TB_KEY"); "" == tb_key {
		log.Panic("need to send telegram bot key in TB_KEY")
	}
	var err error
	bot, err = tgbotapi.NewBotAPI(tb_key)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false


	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		go handleMessage(update.Message)
	}
}
func sendBotMessage(message_id int64, message string) {
	msg := tgbotapi.NewMessage(message_id, message)
	bot.Send(msg)
}

func handleMessage(message *tgbotapi.Message){
	if message == nil {
		return
	}

	chatID := message.Chat.ID
	if message.From.UserName == "abonec" {
		switch message.Text {
		case "/power_on":
			turnPowerOn()
			sendBotMessage(chatID, "power on signal was sent")
		case "/status":
			machineOnline := pingMachine()
			if machineOnline {
				sendBotMessage(chatID, "machine online")
			} else {
				sendBotMessage(chatID, "machine offline")
			}
		default:
			sendBotMessage(chatID, "unknown command")
		}
	} else {
		sendBotMessage(chatID, "you are not authorized for doing it")
	}

	log.Printf("[%s] %s", message.From.UserName, message.Text)
}