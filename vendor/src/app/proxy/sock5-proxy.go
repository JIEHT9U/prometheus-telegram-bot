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
		return nil, fmt.Errorf("Error create new sock5 dialer")
	}

	tranport := &http.Transport{
		Dial: dialer.Dial,
	}
	// setup a http client
	httpClient := &http.Client{
		Transport: tranport,
	}
	return httpClient, nil
}
