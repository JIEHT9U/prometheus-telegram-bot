package proxy

import (
	"fmt"
	"net/http"

	"golang.org/x/net/proxy"
)

func New(network, address string, auth *proxy.Auth, forward proxy.Dialer) (*http.Client, error) {

	dialer, err := proxy.SOCKS5(network, address, auth, forward)
	if err != nil {
		return nil, fmt.Errorf("Error create new")
	}
	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	httpTransport.Dial = dialer.Dial
	return httpClient, nil
}

//socks5://example.com:8080
// func SOCKS5(network, address string, auth *Auth, forward Dialer) (Dialer, error) {
// 	d := socks.NewDialer(network, address)
// 	if forward != nil {
// 		d.ProxyDial = func(_ context.Context, network string, address string) (net.Conn, error) {
// 			return forward.Dial(network, address)
// 		}
// 	}
// 	if auth != nil {
// 		up := socks.UsernamePassword{
// 			Username: auth.User,
// 			Password: auth.Password,
// 		}
// 		d.AuthMethods = []socks.AuthMethod{
// 			socks.AuthMethodNotRequired,
// 			socks.AuthMethodUsernamePassword,
// 		}
// 		d.Authenticate = up.Authenticate
// 	}
// 	return d, nil
// }

/*
func ProxyAwareHttpClient() *http.Client {
	// sane default
	var dialer proxy.Dialer
	// eh, I want the type to be proxy.Dialer but assigning proxy.Direct makes the type proxy.direct
	dialer = proxy.Direct
	proxyServer, isSet := os.LookupEnv("HTTP_PROXY")
	if isSet {
		proxyUrl, err := url.Parse(proxyServer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid proxy url %q\n", proxyUrl)
		}
		dialer, err = proxy.FromURL(proxyUrl, proxy.Direct)
	}

	// setup a http client
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	httpTransport.Dial = dialer.Dial
	return httpClient
}
*/
