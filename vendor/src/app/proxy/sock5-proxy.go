package proxy

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/proxy"
)

func New(network, address string, auth *proxy.Auth, proxyTimeOut time.Duration) (*http.Client, error) {

	dialer, err := proxy.SOCKS5(network, address, auth, proxyTimeOut)

	if err != nil {
		return nil, fmt.Errorf("Error create new")
	}

	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	httpTransport.Dial = dialer.Dial
	return httpClient, nil
}
