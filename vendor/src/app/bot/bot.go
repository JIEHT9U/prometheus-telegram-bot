package bot

import (
	"app/logger"
	"app/options"
	"app/storage"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

type Config struct {
	Token         string
	Logger        *logger.Logger
	clinet        *http.Client
	httpTransport *http.Transport
}

func New(o *options.ServerRunOptions, l *logger.Logger) (*Config, error) {

	var cfg = &Config{
		Logger: l,
		Token:  o.TelegramToken,
	}

	if o.ProxyURL == "" {
		cfg.httpTransport = &http.Transport{}
		cfg.clinet = &http.Client{
			Transport: cfg.httpTransport,
			Timeout:   o.ProxyTimeOut,
		}
		return cfg, nil
	}
	proxyURL, err := url.Parse(o.ProxyURL)
	if err != nil {
		return cfg, fmt.Errorf("Error parsing Tor proxy URL:%s [%s]", o.ProxyURL, err)
	}
	l.InfoEntry().Infof("Use proxy %s", proxyURL.String())

	cfg.httpTransport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	cfg.clinet = &http.Client{
		Transport: cfg.httpTransport,
		Timeout:   o.ProxyTimeOut,
	}

	return cfg, nil
}

func (c *Config) requestGenerator(methods string, params ...string) (*http.Request, error) {

	req, err := http.NewRequest("GET",
		fmt.Sprintf(
			"https://api.telegram.org/bot%s/%s",
			c.Token,
			methods,
		),
		nil)
	if err != nil {
		return nil, fmt.Errorf("Error making GET request [%s]", err)
	}
	return req, nil
}

func (c *Config) GetInfo() (User, error) {
	var user User

	req, err := c.requestGenerator("getMe")
	if err != nil {
		return user, err
	}
	apiResponce, err := c.request(req)
	if err != nil {
		return user, err
	}

	if err := json.Unmarshal(apiResponce.Result, &user); err != nil {
		return user, err
	}

	return user, nil
}

func (c *Config) request(req *http.Request) (APIResponse, error) {
	var api APIResponse
	resp, err := c.clinet.Do(req)
	if err != nil {
		return api, fmt.Errorf("Error request(%s)", req.URL.String())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return api, fmt.Errorf("Responce status code %d != 200", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&api); err != nil {
		return api, fmt.Errorf("Error decode telegram bot APIResponse [%s]", err)
	}

	return api, nil
}

func (c *Config) Send(chatID int64, msg string) (Message, error) {
	botMsg := CreateMessage(chatID, msg)
	botMsg.ParseMode = tgbotapi.ModeHTML
	botMsg.DisableWebPagePreview = true
	return c.send(botMsg)
}

func (c *Config) send(config Chattable) (Message, error) {
	var msg Message
	valueUrl, err := config.values()
	if err != nil {
		return msg, err
	}

	/* 	message, err := makeMessageRequest(config.method(), v)
	   	if err != nil {
	   		return msg, err
	   	} */

	req, err := c.requestGenerator(config.method())
	if err != nil {
		return user, err
	}

	return message, nil
}

/* func (c *Config) makeMessageRequest(endpoint string, params url.Values) (Message, error) {

	resp, err := c.request(endpoint, params)
	if err != nil {
		return Message{}, err
	}

	var message Message
	json.Unmarshal(resp.Result, &message)

	bot.debugLog(endpoint, params, message)

	return message, nil
}
*/
func CreateMessage(chatID int64, text string) MessageConfig {
	return MessageConfig{
		BaseChat: BaseChat{
			ChatID:           chatID,
			ReplyToMessageID: 0,
		},
		Text: text,
		DisableWebPagePreview: false,
	}
}

// type BOT_COMMAND string
type TelegramBot struct {
	Bot     *tgbotapi.BotAPI
	Storage storage.Storage
	logger  *logger.Logger
}

func chechProxyRequire(o *options.ServerRunOptions) bool {
	if o.ProxyURL != "" {
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
	var client *http.Client

	if chechProxyRequire(o) {
		l.InfoEntry().Infof("Proxy %s connections...", o.ProxyURL)

		proxyURL, err := url.Parse(o.ProxyURL)
		if err != nil {
			return tBot, fmt.Errorf("Error parse proxy url [%s]", err)
		}

		/* 		switch proxyURL.Scheme {
		   		case "socks5":
		   			password, _ := proxyURL.User.Password()
		   			client, err = proxy.New("tcp", strings.Join([]string{proxyURL.Hostname(), proxyURL.Port()}, ":"), &p.Auth{User: proxyURL.User.Username(), Password: password}, o.ProxyTimeOut)
		   			if err != nil {
		   				return tBot, err
		   			}
		   		default:
		   			proxyTransport := &http.Transport{
		   				Proxy: http.ProxyURL(proxyURL),
		   			}
		   			client = &http.Client{
		   				Transport: proxyTransport,
		   				// Timeout:   o.ProxyTimeOut,

		   			}
		   		} */

		// password, _ := proxyURL.User.Password()
		// client, err = proxy.New("tcp", "95.85.30.96:1090", &p.Auth{User: "golang", Password: "golang777"}, 0)
		// client, err = proxy.New("tcp4", "socks5://95.85.30.96:1090", &p.Auth{User: proxyURL.User.Username(), Password: password}, o.ProxyTimeOut)
		proxyTransport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
		client = &http.Client{
			Transport: proxyTransport,
			Timeout:   o.ProxyTimeOut,
		}
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
	l.InfoEntry().Info("Success bot connections")

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
