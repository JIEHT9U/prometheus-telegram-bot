package bot

import (
	"app/logger"
	"app/options"
	"app/proxy"
	"app/storage"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	p "golang.org/x/net/proxy"
)

// type BOT_COMMAND string
type TelegramBot struct {
	Bot     *tgbotapi.BotAPI
	Storage storage.Storage
	logger  *logger.Logger
}

func chechProxyRequire(o *options.ServerRunOptions) bool {
	if o.ProxyPassword != "" && o.ProxyUser != "" && o.ProxyNetwork != "" {
		return true
	}
	return false
}

/* func defaulDialer() *http.Client {

	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	httpTransport.Dial = net.Dialer{}
	return httpClient, nil
} */

//Create ...
func Create(o *options.ServerRunOptions, l *logger.Logger) (TelegramBot, error) {

	var tBot TelegramBot
	var err error
	var bot *tgbotapi.BotAPI

	if chechProxyRequire(o) {
		l.InfoEntry().Infof("Proxy %s connections...", o.ProxyURL)

		client, err := proxy.New(o.ProxyNetwork, o.ProxyURL, &p.Auth{User: o.ProxyUser, Password: o.ProxyPassword}, o.ProxyTimeOut)
		if err != nil {
			return tBot, err
		}
		bot, err = tgbotapi.NewBotAPIWithClient(o.TelegramToken, client)
	} else {
		bot, err = tgbotapi.NewBotAPI(o.TelegramToken)
	}

	if err != nil {
		return tBot, errors.Wrap(err, "Telegram bot")
	}
	// l.InfoEntry().Info("Success bot connections")

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
