// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proxy

import (
	"time"

	"golang.org/x/net/internal/socks"
)

// SOCKS5 returns a Dialer that makes SOCKSv5 connections to the given
// address with an optional username and password.
// See RFC 1928 and RFC 1929.
func SOCKS5(network, address string, auth *Auth, proxyTimeOut time.Duration) (Dialer, error) {

	dialer := socks.NewDialer(network, address, proxyTimeOut)

	if auth != nil {
		usernamePassword := socks.UsernamePassword{
			Username: auth.User,
			Password: auth.Password,
		}
		dialer.AuthMethods = []socks.AuthMethod{
			socks.AuthMethodNotRequired,
			socks.AuthMethodUsernamePassword,
		}
		dialer.Authenticate = usernamePassword.Authenticate
	}
	return dialer, nil
}
