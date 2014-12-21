package netutil

import (
	"net"
	"time"
)

type DialService interface {
	Listen(net, laddr string) (net.Listener, error)
	Dial(network, address string) (net.Conn, error)
	DialTimeout(network, address string, timeout time.Duration) (net.Conn, error)
}

var (
	gDialServices map[string]DialService
)

func AddDialService(n string, s DialService) {
	if gDialServices == nil {
		gDialServices = make(map[string]DialService)
	}
	gDialServices[n] = s
}

func GetDialService(n string) DialService {
	return gDialServices[n]
}

func Dial(network, address string) (net.Conn, error) {
	s := GetDialService(network)
	if s == nil {
		return net.Dial(network, address)
	}
	return s.Dial(network, address)
}

func DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	s := GetDialService(network)
	if s == nil {
		return net.DialTimeout(network, address, timeout)
	}
	return s.DialTimeout(network, address, timeout)
}

func Listen(network, laddr string) (net.Listener, error) {
	s := GetDialService(network)
	if s == nil {
		return net.Listen(network, laddr)
	}
	return s.Listen(network, laddr)
}
