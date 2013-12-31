package httputil

import (
	"net"
	"net/http"
	"time"
)

func TimeoutDialer(timeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	now := time.Now()
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, timeout)
		if err != nil {
			return nil, err
		}
		tm := now.Add(timeout)
		conn.SetDeadline(tm)
		return conn, nil
	}
}

func NewHttpClient(timeout time.Duration) *http.Client {

	return &http.Client{
		Transport: &http.Transport{
			Dial: TimeoutDialer(timeout),
		},
	}
}
