// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package socks

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strconv"
	"time"
)

var (
	noDeadline   = time.Time{}
	aLongTimeAgo = time.Unix(1, 0)
)

func (d *Dialer) upgrader(c net.Conn, address string) (net.Addr, error) {
	host, port, err := splitHostPort(address)
	if err != nil {
		return nil, err
	}

	authMsg, err := d.authorizationMsg()
	if err != nil {
		return nil, fmt.Errorf("Error create socks5 auth Msg [%s]", err)
	}

	if _, err := c.Write(authMsg); err != nil {
		return nil, fmt.Errorf("Error write socks5 auth Msg in connection [%s]", err)
	}

	if err := d.validateSocksVersion(c); err != nil {
		return nil, fmt.Errorf("Error validate socks5 version [%s]", err)
	}

	cmdByte, err := d.useCmd(host, port)
	if err != nil {
		return nil, fmt.Errorf("Error use Cmd socks5 [%s]", err)
	}

	if _, err := c.Write(cmdByte); err != nil {
		return nil, fmt.Errorf("Error write socks5 cmd Byte in connection [%s]", err)
	}

	return validateResponce(c)
}

func validateResponce(c net.Conn) (net.Addr, error) {
	var res = make([]byte, 4)
	var l = 2
	a := &Addr{}

	if _, err := io.ReadFull(c, res); err != nil {
		return a, fmt.Errorf("validate comand responce [%s]", err)
	}

	if res[0] != Version5 {
		return a, errors.New("unexpected protocol version " + strconv.Itoa(int(res[0])))
	}

	if err := Reply(res[1]); err != StatusSucceeded {
		return a, errors.New("unknown error " + err.String())
	}

	if res[2] != 0 {
		return a, errors.New("non-zero reserved field")
	}

	switch res[3] {
	case AddrTypeIPv4:
		l += net.IPv4len
		a.IP = make(net.IP, net.IPv4len)
	case AddrTypeIPv6:
		l += net.IPv6len
		a.IP = make(net.IP, net.IPv6len)
	case AddrTypeFQDN:
		r, err := ioutil.ReadAll(c)
		if err != nil {
			return a, err
		}
		l += int(r[0])
	default:
		return a, errors.New("unknown address type " + strconv.Itoa(int(res[3])))
	}

	b := make([]byte, l)
	if _, err := io.ReadFull(c, b); err != nil {
		return a, err
	}

	if a.IP != nil {
		copy(a.IP, b)
	} else {
		a.Name = string(b[:len(b)-2])
	}
	a.Port = int(b[len(b)-2])<<8 | int(b[len(b)-1])
	return a, nil
}

func (d *Dialer) useCmd(host string, port int) ([]byte, error) {
	var cmd bytes.Buffer

	cmd.WriteByte(Version5)
	cmd.WriteByte(byte(d.cmd))
	cmd.WriteByte(0)

	if ip := net.ParseIP(host); ip != nil {
		if ip4 := ip.To4(); ip4 != nil {
			cmd.WriteByte(AddrTypeIPv4)
			cmd.WriteByte(AddrTypeIPv4)
			for _, ipBytes := range ip4 {
				cmd.WriteByte(ipBytes)
			}

		} else if ip6 := ip.To16(); ip6 != nil {
			cmd.WriteByte(AddrTypeIPv6)
			for _, ipBytes := range ip6 {
				cmd.WriteByte(ipBytes)
			}

		} else {
			return nil, errors.New("unknown address type")
		}
	} else {
		hostLen := len(host)
		if hostLen > 255 {
			return nil, errors.New("FQDN too long")
		}
		cmd.WriteByte(AddrTypeFQDN)
		cmd.WriteByte(byte(hostLen))
		for _, hBytes := range []byte(host) {
			cmd.WriteByte(hBytes)
		}
	}

	cmd.WriteByte(byte(port >> 8))
	cmd.WriteByte(byte(port))

	return cmd.Bytes(), nil
}

func (d *Dialer) validateSocksVersion(c net.Conn) error {
	res := make([]byte, 2)

	if _, err := io.ReadFull(c, res); err != nil {
		return fmt.Errorf("validate socks version [%s]", err)
	}

	if res[0] != Version5 {
		return errors.New("unexpected protocol version " + strconv.Itoa(int(res[0])))
	}

	am := AuthMethod(res[1])
	if am == AuthMethodNoAcceptableMethods {
		return errors.New("no acceptable authentication methods")
	}

	if err := d.Authenticate(c, am); err != nil {
		return err
	}

	return nil
}

func (d *Dialer) authorizationMsg() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteByte(Version5)
	if len(d.AuthMethods) == 0 || d.Authenticate == nil {
		buf.WriteByte(1)
		buf.WriteByte(byte(AuthMethodNotRequired))
	} else {
		aMethods := d.AuthMethods
		lenAMethods := len(aMethods)
		if lenAMethods > 255 {
			return nil, errors.New("too many authentication methods")
		}
		buf.WriteByte(byte(lenAMethods))
		for _, am := range aMethods {
			buf.WriteByte(byte(am))
		}
	}

	return buf.Bytes(), nil
}

func splitHostPort(address string) (string, int, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return "", 0, err
	}
	portnum, err := strconv.Atoi(port)
	if err != nil {
		return "", 0, err
	}
	if 1 > portnum || portnum > 0xffff {
		return "", 0, errors.New("port number out of range " + port)
	}
	return host, portnum, nil
}
