package main

import (
	"app"
	"app/logger"
	"app/options"
	"fmt"

	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/pflag"
)

// var onlyOneSignalHandler = make(chan struct{})
// var shutdownSignals = []os.Signal{os.Interrupt}
// Version is the tagged version or "dev"

var Version = "dev"

// BuildDate is the date of the release build or ""
var BuildDate = ""

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}

func SetupSignalHandler() (stopCh <-chan struct{}) {
	// close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}

// WordSepNormalizeFunc changes all flags that contain "_" separators
func WordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.Replace(name, "_", "-", -1))
	}
	return pflag.NormalizedName(name)
}

func InitFlags() {
	pflag.CommandLine.SetNormalizeFunc(WordSepNormalizeFunc)
	pflag.CommandLine.SortFlags = false
	pflag.Parse()
}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	s := options.NewServerRunOptions()

	s.AddFlags(pflag.CommandLine)
	InitFlags()

	l, err := logger.New(s.Debug, s.JSON)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	l.InfoEntry().Debug("Enable Debug mode")

	if BuildDate == "" {
		BuildDate = time.Now().Format("15:04:05 Jan 2, 2006")
	}

	l.InfoEntry().Infof("Prometheus telegram bot version %s built %v\n", Version, BuildDate)

	stopCh := SetupSignalHandler()
	if err := app.Run(s, stopCh, l); err != nil {
		l.ErrEntry().Error(err)
		os.Exit(1)
	}

}
