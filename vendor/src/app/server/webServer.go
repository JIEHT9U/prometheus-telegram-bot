package server

import (
	"app/bot"
	"app/logger"
	"app/msg"
	"app/options"
	t "app/template"
	"context"
	"encoding/json"
	"fmt"
	tmplhtml "html/template"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WebServerConfig struct {
	Addr     string
	Logger   *logger.Logger
	Debug    bool
	Shutdown <-chan struct{}
	MsgSize  int
}

func NewWebConfig(s *options.ServerRunOptions, l *logger.Logger, shutdown <-chan struct{}) WebServerConfig {
	return WebServerConfig{
		Addr:     strings.Join([]string{s.BindAddress.String(), ":", strconv.Itoa(s.BindPort)}, ""),
		Debug:    s.Debug,
		Logger:   l,
		Shutdown: shutdown,
		MsgSize:  s.MessageSizeBytes,
	}
}

func (c WebServerConfig) CreateWebServer(handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    c.Addr,
		Handler: handler,
	}
}

type wsMsg struct {
	BotFormat string
	Received  string
	Count     int
}

func (c WebServerConfig) GetHandler(bot bot.TelegramBot, tmps map[string]*tmplhtml.Template) http.Handler {
	r := mux.NewRouter()

	upgrader := websocket.Upgrader{
		ReadBufferSize:    4096,
		WriteBufferSize:   4096,
		EnableCompression: true,
		Subprotocols:      []string{"websocket"},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	wsChanel := make(chan wsMsg, 1)

	r.HandleFunc("/ws/{ws_id}", func(w http.ResponseWriter, r *http.Request) {
		c.Logger.InfoEntry().Infof("New Ws connection:%s", r.URL)
		// serveWs(hub, w, r)
		defer c.Logger.InfoEntry().Infof("Defer Ws connection:%s", r.URL)

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			c.Logger.ReqError(w, err)
			c.Logger.ErrEntry().Error(err)
			return
		}
		defer ws.Close()

		go func() {

			for {
				_, msg, err := ws.ReadMessage()
				if err != nil {
					c.Logger.ErrEntry().Error(err)
					break
				}
				c.Logger.InfoEntry().Infof("msg:%s", string(msg))
			}
		}()

		for {
			select {
			case wsMsg := <-wsChanel:
				err := ws.WriteJSON(wsMsg.Received)
				if err != nil {
					c.Logger.ErrEntry().Errorf("Error send msg in WebSocket:%s", r.URL)
					break
				}
				// fmt.Println(wsMsg)
			}
		}

	})

	r.HandleFunc("/alert/{template_name}/{chat_id}", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		alertsMsg, receiveMsg, err := msg.Parser(r.Body)
		if err != nil {
			c.Logger.ReqError(w, err)
			c.Logger.ErrEntry().Error(err)
			c.Logger.ErrEntry().Errorf("Received post request:\n %s", receiveMsg)
			return
		}

		// c.Logger.InfoEntry().Debugf("Received post request:\n %s", receiveMsg)
		// c.Logger.InfoEntry().Debug("+-----------------------------------------------------------+")
		// c.Logger.InfoEntry().Debug("+-----------------------------------------------------------+")
		// c.Logger.InfoEntry().Debugf("Alert Msg:\n %s", printStrucPretty(alertsMsg))

		templateName, chatID, err := RetrieveTemplateNameAndChatID(r)

		if err != nil {
			c.Logger.ReqError(w, err)
			c.Logger.ErrEntry().Error(err)
			return
		}

		c.Logger.InfoCT(templateName, chatID, "Received prometheus alert")

		template, err := t.Find(tmps, templateName)

		msg, ok := err.(*t.ErrDefaultTempatestruct)
		if err != nil && !ok {
			c.Logger.ReqError(w, err)
			c.Logger.ErrorCT(templateName, chatID, err)
			return
		}

		if ok {
			c.Logger.InfoEntry().Info(msg)
			templateName = t.DEFAULT_TEMPLATE_NAME
		}

		finalMSg, err := t.ExecuteTextString(template, alertsMsg)
		if err != nil {
			err := fmt.Errorf("Error exec %s templates [%s]", templateName, err)
			c.Logger.ReqError(w, err)
			c.Logger.ErrorCT(templateName, chatID, err)
			return
		}

		c.Logger.InfoEntry().Debugf("Final MSg :\n %s", finalMSg)

		select {
		case wsChanel <- wsMsg{
			BotFormat: finalMSg,
			Count:     c.MsgSize,
			Received:  receiveMsg,
		}:
			c.Logger.InfoEntry().Info("Success send msg in wsMsg chanel")
		case <-time.After(time.Microsecond * 500):
			c.Logger.ErrEntry().Error("Error send msg in wsMsg chanel timeout 500ms")
		}

		if len(finalMSg) > c.MsgSize {
			for i, ss := range SplitMsg(finalMSg, c.MsgSize) {
				_, err = bot.Send(chatID, ss)
				if err != nil {
					err := fmt.Errorf("Error send  msg telegram bot: %s", err)
					c.Logger.ReqError(w, err)
					c.Logger.ErrorCT(templateName, chatID, err)
					return
				}
				c.Logger.InfoCT(templateName, chatID, "Succes send bot msg №", i+1)
			}

			c.Logger.InfoCT(templateName, chatID, "Succes send all bot msg ")
			return
		}

		_, err = bot.Send(chatID, finalMSg)
		if err != nil {
			err := fmt.Errorf("Error send  msg telegram bot: %s", err)
			c.Logger.ReqError(w, err)
			c.Logger.ErrorCT(templateName, chatID, err)
			return
		}
		c.Logger.InfoCT(templateName, chatID, "Succes send bot msg")
		// l.DebugfCT(templateName, chatID, "Succes send bot msg: %s", printStrucPretty(resultSendMsg))
	})
	return r
}

func SplitMsg(s string, n int) []string {

	var concatString string
	var floor = int(math.Floor(float64(len(s)) / float64(n)))
	var res = make([]string, 0, floor)

	for _, s := range strings.Split(s, "\n") {
		fmtString := fmt.Sprintf("%s \n", s)
		if (len(concatString) + len(fmtString)) <= n {
			concatString += fmtString
			continue
		}
		res = append(res, concatString)
		concatString = ""
	}

	if len(res) != floor+1 {
		res = append(res, concatString)
	}

	return res
}

//Run func will be run web server
func (c WebServerConfig) Run(server *http.Server) {

	go func() {
		c.Logger.InfoEntry().Infof("Run web server on: %v", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			c.Logger.ErrEntry().Fatalf("Err web server [ %s ]", err)
		}
	}()

	<-c.Shutdown

	c.Logger.InfoEntry().Info("Run gracefull  web server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	server.Shutdown(ctx)

	<-ctx.Done()
	c.Logger.InfoEntry().Info("Graceful shutdown web server")
}

//RetrieveTemplateNameAndChatID return (templateName , chatID , error )
func RetrieveTemplateNameAndChatID(r *http.Request) (string, int64, error) {

	vars := mux.Vars(r)

	templateName, ok := vars["template_name"]
	if !ok {
		return "", 0, fmt.Errorf("Not found template_name in request %s", r.URL.String())
	}

	chatID, ok := vars["chat_id"]
	if !ok {
		return "", 0, fmt.Errorf("Not found chat_id in request %s", r.URL.String())
	}

	chatIDInt64, err := strconv.ParseInt(chatID, 0, 64)
	if err != nil {
		return "", 0, fmt.Errorf("Error parse chat_id: %s in int64. [%s]", chatID, err)
	}
	return templateName, chatIDInt64, nil
}

func printStrucPretty(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
