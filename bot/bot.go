package bot

import (
	"fmt"
	"log"

	"github.com/webbsalad/weather-bot/weather"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api *tgbot.BotAPI
}

func NewBot(token string) *Bot {
	bot, err := tgbot.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	return &Bot{api: bot}
}

func (b *Bot) Start() {
	u := tgbot.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				b.handleCommand(update.Message)
			} else if update.Message.Text != "" {
				b.handleText(update.Message)
			}
		}
	}
}

func (b *Bot) handleCommand(message *tgbot.Message) {

	switch message.Command() {
	case "start":
		buttons := tgbot.NewReplyKeyboard(
			tgbot.NewKeyboardButtonRow(
				tgbot.NewKeyboardButton("Москва"),
				tgbot.NewKeyboardButton("Санкт-Петербург"),
				tgbot.NewKeyboardButton("Улан-Удэ"),
			),
		)
		msg := tgbot.NewMessage(message.Chat.ID, "Выберите город, чтобы узнать погоду:")
		msg.ReplyMarkup = buttons
		b.api.Send(msg)
	default:
		msg := tgbot.NewMessage(message.Chat.ID, "Не знаю такой команды")
		b.api.Send(msg)
	}
}

func (b *Bot) handleText(message *tgbot.Message) {
	city := message.Text
	weatherData, err := weather.Get(city)
	if err != nil {
		msg := tgbot.NewMessage(message.Chat.ID, "Извините, не удалось получить погоду.")
		b.api.Send(msg)
		return
	}

	if len(weatherData.Weather) == 0 {
		msg := tgbot.NewMessage(message.Chat.ID, "Извините, нет информации о погоде.")
		b.api.Send(msg)
		return
	}

	response := fmt.Sprintf("%s:\nТемпература: %.2f°C\nОписание: %s", city, weatherData.Main.Temp-273.15, weatherData.Weather[0].Description)
	msg := tgbot.NewMessage(message.Chat.ID, response)
	b.api.Send(msg)
}
