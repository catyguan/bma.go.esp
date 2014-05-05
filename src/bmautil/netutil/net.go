package netutil

import (
	"net"
	"strings"
	"sync"
)

type ConnGroup struct {
	conns map[net.Conn]bool
	lock  sync.Mutex
}

func NewConnGroup() *ConnGroup {
	r := new(ConnGroup)
	r.conns = make(map[net.Conn]bool)
	return r
}

func (this *ConnGroup) Add(c net.Conn) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	if _, ok := this.conns[c]; ok {
		return false
	}
	this.conns[c] = true
	return true
}

func (this *ConnGroup) Remove(c net.Conn) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	if _, ok := this.conns[c]; ok {
		delete(this.conns, c)
		return true
	}
	return false
}

func (this *ConnGroup) CloseAll() {
	this.lock.Lock()
	defer this.lock.Unlock()
	for conn, _ := range this.conns {
		delete(this.conns, conn)
		conn.Close()
	}
}

type ChannelCloseListener func()

type Channel struct {
	net.Conn
	listeners  []ChannelCloseListener
	Properties map[string]interface{}
}

func NewChannel(c net.Conn) *Channel {
	r := new(Channel)
	r.Conn = c
	r.Properties = make(map[string]interface{})
	return r
}

func notify(lis ChannelCloseListener) {
	defer func() {
		recover()
	}()
	lis()
}

func (this *Channel) CloseChannel() {
	this.Close()
	if this.listeners != nil {
		for _, lis := range this.listeners {
			notify(lis)
		}
		this.listeners = make([]ChannelCloseListener, 0)
	}
}

func (this *Channel) AddListener(f ChannelCloseListener) {
	if this.listeners == nil {
		this.listeners = []ChannelCloseListener{f}
	} else {
		this.listeners = append(this.listeners, f)
	}
}

func IpMatch(addr net.IP, list []string) (bool, string) {
	if list == nil {
		return false, ""
	}
	for _, s := range list {
		v := strings.TrimSpace(s)
		_, ipNet, err := net.ParseCIDR(v)
		if err != nil {
			ipHost := net.ParseIP(v)
			if ipHost != nil {
				if ipHost.Equal(addr) {
					return true, v
				}
			}
		} else {
			if ipNet.Contains(addr) {
				return true, v
			}
		}
	}
	return false, ""
}

func IpAccept(address string, whiteList []string, blackList []string, notMatchReturn bool) (bool, string) {
	addr := strings.Split(address, ":")
	raddr := net.ParseIP(addr[0])
	if ok, msg := IpMatch(raddr, blackList); ok {
		return false, "BLACK: " + msg
	}
	if whiteList != nil && len(whiteList) > 0 {
		if ok, msg := IpMatch(raddr, whiteList); ok {
			return true, "WHITE: " + msg
		}
		return false, "NOT WHITE: " + address
	}
	return notMatchReturn, ""
}
