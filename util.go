package main

import (
	"net"
	"net/http"
	"regexp"
	"time"
)

var reg *regexp.Regexp = regexp.MustCompile("[^a-zA-Z0-9-]+")

const MAX_RETRY = 20
const WAIT_RETRY = time.Second * 1

func getFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

func filterInvalidHeaders(headers http.Header) {
	for header := range headers {
		if reg.MatchString(header) {
			delete(headers, header)
		}
	}
}
