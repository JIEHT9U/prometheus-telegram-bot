package app

import (
	"app/bot"
	"app/logger"
	m "app/mapping"
	"app/options"
	"app/server"
	t "app/template"
	"app/utilerrors"
	"fmt"
)

func Run(runOptions *options.ServerRunOptions, stopCh <-chan struct{}, l *logger.Logger) error {

	// To help debugging, immediately log version

	if errs := runOptions.Validate(); len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}

	mapsInstance, err := m.GetNames(runOptions.MappingNamePaths)
	if err != nil {
		return err
	}

	tmps, err := t.Load(mapsInstance, runOptions.TemplatePaths)
	if err != nil {
		return err
	}

	if len(tmps) <= 0 {
		return fmt.Errorf("Not fount templates in %s", t.AlignmentPath(runOptions.TemplatePaths))
	}

	for name := range tmps {
		l.InfoEntry().Infof("Success read template: %s ", name)
	}

	b, err := bot.Create(runOptions, l)
	if err != nil {
		return fmt.Errorf("Err create telegram bot:[%s]", err)
	}
	if b.Bot == nil {
		return fmt.Errorf("Err proxy connection")
	}
	l.InfoEntry().Infof("Authorized on account [ %s ]", b.Bot.Self.UserName)

	err, _ = b.Run(stopCh)
	if err != nil {
		return err
	}

	webCfg := server.NewWebConfig(runOptions, l, stopCh)
	webCfg.Run(webCfg.CreateWebServer(webCfg.GetHandler(b, tmps)))

	return nil
}
