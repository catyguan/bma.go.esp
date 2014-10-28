package httputil

import (
	"net"
	"net/http"
	"strings"
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

func IsMultipartForm(r *http.Request) bool {
	ct := r.Header.Get("Content-Type")
	return strings.HasPrefix(ct, "multipart/form-data")
}

func Prepare(r *http.Request, maxMemory int64) error {
	if IsMultipartForm(r) {
		return r.ParseMultipartForm(maxMemory)
	} else {
		return r.ParseForm()
	}
}
