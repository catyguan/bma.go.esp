package netutil

import (
	"net"
	"strings"
)

type commandConnGroup struct {
	cmd   int // 0-closeAll 1-add 2-remove 3-removeAndClose
	conn  net.Conn
	event func()
}
type ConnGroup struct {
	conns   map[net.Conn]bool
	command chan commandConnGroup
}

func NewConnGroup(size int) *ConnGroup {
	r := new(ConnGroup)
	r.command = make(chan commandConnGroup, size)
	r.conns = make(map[net.Conn]bool)
	go r.run()
	return r
}

func (this *ConnGroup) Add(c net.Conn, cb func()) (r bool) {
	defer func() {
		err := recover()
		if err != nil {
			r = false
		}
	}()
	this.command <- commandConnGroup{1, c, cb}
	r = true
	return
}

func (this *ConnGroup) Remove(c net.Conn, cb func()) (r bool) {
	defer func() {
		err := recover()
		if err != nil {
			r = false
		}
	}()
	this.command <- commandConnGroup{2, c, cb}
	r = true
	return
}

func (this *ConnGroup) RemoveAll(cb func()) (r bool) {
	defer func() {
		err := recover()
		if err != nil {
			r = false
		}
	}()
	this.command <- commandConnGroup{0, nil, cb}
	r = true
	return
}

func (this *ConnGroup) Close() {
	close(this.command)
}

func callback(cb func()) {
	if cb != nil {
		defer func() {
			recover()
		}()
		cb()
	}
}

func (this *ConnGroup) run() {
	defer this.closeAll()
	for cobj := range this.command {
		switch cobj.cmd {
		case 0:
			this.closeAll()
		case 1:
			this.conns[cobj.conn] = true
		case 2:
			if _, ok := this.conns[cobj.conn]; ok {
				cobj.conn.Close()
				delete(this.conns, cobj.conn)
			}
		}
		if cobj.event != nil {
			callback(cobj.event)
		}
	}
}

func (this *ConnGroup) closeAll() {
	for conn, _ := range this.conns {
		conn.Close()
	}
	this.conns = make(map[net.Conn]bool)
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
