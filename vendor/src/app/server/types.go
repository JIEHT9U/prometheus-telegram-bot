package server

import (
	"app/options"
	"net/http"
	"time"
)

type PrometheusBotServer struct {
	webServer       *http.Server
	Shutdown        <-chan struct{}
	ShutdownTimeout time.Duration
}

//DefaultServerCfg return default conf
func DefaultCfg(o *options.ServerRunOptions, stopCh <-chan struct{}) *PrometheusBotServer {
	return &PrometheusBotServer{
		ShutdownTimeout: time.Second * 10,
		Shutdown:        stopCh,
	}
}

/*
//SetWebServerConfig func
func (s *GenericSmartIDServer) SetWebServerConfig(server *http.Server) {
	s.webServer = server
}
*/
