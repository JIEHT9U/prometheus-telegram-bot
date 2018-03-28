package bot

import (
	"app/logger"
	"app/storage"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

// type BOT_COMMAND string
type TelegramBot struct {
	Bot     *tgbotapi.BotAPI
	Storage storage.Storage
	logger  *logger.Logger
}

//Create(token string) (<-chan tgbotapi.Update, *tgbotapi.BotAPI, error) {
func Create(token string, l *logger.Logger) (TelegramBot, error) {
	var tBot TelegramBot
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return tBot, errors.Wrap(err, "Telegram bot")
	}

	return TelegramBot{Bot: bot, logger: l}, nil
}

func (tb TelegramBot) getUpdateChanel() (<-chan tgbotapi.Update, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := tb.Bot.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}
	return updates, nil
}

func (tb TelegramBot) Send(chatID int64, msg string) (tgbotapi.Message, error) {
	botMsg := tgbotapi.NewMessage(chatID, msg)
	botMsg.ParseMode = tgbotapi.ModeHTML
	botMsg.DisableWebPagePreview = true
	return tb.Bot.Send(botMsg)
}

func (tb TelegramBot) Run(shutdown <-chan struct{}) (error, chan interface{}) {
	tb.logger.InfoEntry().Info("Success start received telegram msg")
	recives := make(chan interface{}, 1)

	updates, err := tb.getUpdateChanel()
	if err != nil {
		return err, nil
	}
	go func() {
		defer tb.logger.InfoEntry().Info("Shutdown telegram bot")
		for {
			select {
			case update := <-updates:

				introduce := func() {
					chatID := update.Message.Chat.ID
					_, err := tb.Send(chatID, fmt.Sprintf(`ℹ <i>Chat id is</i> <b> '%d' </b>`, chatID))

					if err != nil {
						tb.logger.ErrEntry().WithField("chat_id", chatID).Error(errors.Wrap(err, "Error send  chat_id when append or remove telegram bot in group"))
					}

					tb.logger.InfoEntry().WithField("chat_id", chatID).Info("Success send  chat_id when append or remove telegram bot in group ")
				}

				if update.Message.NewChatMembers != nil && len(*update.Message.NewChatMembers) > 0 {
					for _, member := range *update.Message.NewChatMembers {
						if member.UserName == tb.Bot.Self.UserName && update.Message.Chat.Type == "group" {
							introduce()
						}
					}
				} else if update.Message != nil && update.Message.Text != "" {
					introduce()
				}

			case <-shutdown:

				return
			}
		}
	}()
	return nil, recives
}
