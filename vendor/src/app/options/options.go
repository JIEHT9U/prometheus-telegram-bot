package options

import (
	"fmt"
	"net"
	"time"

	"github.com/spf13/pflag"
)

type ServerRunOptions struct {
	TelegramToken    string
	ProxyUser        string
	ProxyPassword    string
	ProxyURL         string
	ProxyNetwork     string
	ProxyTimeOut     time.Duration
	TemplatePaths    []string
	TimeZone         string
	TimeOutFormat    string
	MessageSizeBytes int
	BindAddress      net.IP
	BindPort         int
	Debug            bool
	JSON             bool
	MappingNamePaths []string
}

// NewServerRunOptions creates a new ServerRunOptions object with default parameters
func NewServerRunOptions() *ServerRunOptions {
	return &ServerRunOptions{
		BindAddress:      net.ParseIP("0.0.0.0"),
		BindPort:         9087,
		TemplatePaths:    []string{"template/*.tmpl"},
		MappingNamePaths: []string{"mapping/*.yaml"},
		MessageSizeBytes: 2048,
		ProxyNetwork:     "tcp",
		ProxyTimeOut:     time.Millisecond * 1500,
	}
}

// AddFlags adds flags for a specific APIServer to the specified FlagSet
func (server *ServerRunOptions) AddFlags(fs *pflag.FlagSet) {

	fs.StringVar(&server.TelegramToken, "telegram-token", server.TelegramToken, ""+
		"Telegram token")

	fs.StringVar(&server.ProxyUser, "proxy-user", server.ProxyUser, ""+
		"Proxy user name")

	fs.StringVar(&server.ProxyPassword, "proxy-password", server.ProxyPassword, ""+
		"Proxy password")

	fs.StringVar(&server.ProxyURL, "proxy-url", server.ProxyURL, ""+
		"Proxy url path")

	fs.StringSliceVar(&server.TemplatePaths, "template-paths", server.TemplatePaths, ""+
		"Template path")

	fs.StringSliceVar(&server.MappingNamePaths, "mapping-files-paths", server.MappingNamePaths, ""+
		"Paths for name mapping files")

	fs.StringVar(&server.TimeZone, "time-zone", server.TimeZone, ""+
		"Time zone")

	fs.StringVar(&server.TimeOutFormat, "time-out-format", server.TimeOutFormat, ""+
		"Time out format")

	fs.IntVar(&server.MessageSizeBytes, "message-size-bytes", server.MessageSizeBytes, ""+
		"Telegam message bytr suze (Current maximum length is 4096 UTF8 characters)")

	fs.IPVar(&server.BindAddress, "bind-address", server.BindAddress, ""+
		"Web bind address. ")

	fs.IntVar(&server.BindPort, "bind-port", server.BindPort, ""+
		"Web bind port.")

	fs.BoolVarP(&server.Debug, "debug", "d", server.Debug, ""+
		"Enable debug mod")

	fs.BoolVarP(&server.JSON, "json", "j", server.JSON, ""+
		"Output log format JSON or Systemd")

	fs.DurationVar(&server.ProxyTimeOut, "proxy-timeout", server.ProxyTimeOut, ""+
		"Proxy time out")
}

func (options *ServerRunOptions) Validate() []error {
	var errors []error

	if err := portValidate(options.BindPort); err != nil {
		errors = append(errors, err)
	}

	if err := telegramTokenValidate(options.TelegramToken); err != nil {
		errors = append(errors, err)
	}

	if err := messageSizeBytesValidation(options.MessageSizeBytes); err != nil {
		errors = append(errors, err)
	}

	return errors

}

func messageSizeBytesValidation(size int) error {
	if size > 4096 {
		return fmt.Errorf("--message-size-bytes  %d must be < 4096", size)
	}
	return nil
}

func portValidate(bindPort int) error {
	if bindPort < 0 {
		return fmt.Errorf("--bind-port %v  должен быть в диапазоне > 0", bindPort)
	}
	return nil
}

func telegramTokenValidate(token string) error {
	if token == "" {
		return fmt.Errorf("--telegram-token %s  should be is not empty", token)
	}
	return nil
}
