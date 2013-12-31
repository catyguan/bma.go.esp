package espnet

import (
	"bmautil/socket"
	"logger"
	"net"
	"time"
)

// simple dial
func Dial(name string, cfg *DialConfig, sinit socket.SocketInit) (*socket.Socket, error) {
	if err := cfg.Valid(); err != nil {
		return nil, err
	}

	var conn net.Conn
	var err error
	if cfg.TimeoutMS > 0 {
		conn, err = net.Dial(cfg.Net, cfg.Address)
	} else {
		conn, err = net.DialTimeout(cfg.Net, cfg.Address, time.Duration(cfg.TimeoutMS)*time.Millisecond)
	}
	if err != nil {
		logger.Debug(tag, "dial (%s %s)fail - %s", cfg.Net, cfg.Address, err.Error())
		return nil, err
	}
	sock := socket.NewSocket(conn, 32, 0)
	err = sock.Start(sinit)
	if err != nil {
		return nil, err
	}
	return sock, nil
}

type SimpleDialPool struct {
	name         string
	config       *DialConfig
	socketInit   socket.SocketInit
	channelCoder string
}

func NewSimpleDialPool(n string, cfg *DialConfig, sinit socket.SocketInit, chcoder string) *SimpleDialPool {
	this := new(SimpleDialPool)
	this.name = n
	this.config = cfg
	this.socketInit = sinit
	this.channelCoder = chcoder
	return this
}

func (this *SimpleDialPool) NewChannel() (Channel, error) {
	sock, err := Dial(this.name, this.config, this.socketInit)
	if err != nil {
		return nil, err
	}
	r := NewSocketChannel(sock, this.channelCoder)
	return r, nil
}
