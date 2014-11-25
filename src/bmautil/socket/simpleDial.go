package socket

import (
	"logger"
	"net"
	"time"
)

// simple dial
func Dial(name string, cfg *DialConfig, sinit SocketInit) (*Socket, error) {
	return Dial2(name, cfg, sinit, true)
}

func Dial2(name string, cfg *DialConfig, sinit SocketInit, log bool) (*Socket, error) {
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
		if log {
			logger.Debug(tag, "dial (%s %s) fail - %s", cfg.Net, cfg.Address, err.Error())
		}
		return nil, err
	}
	sock := NewSocket(conn, 32, 0)
	err = sock.Start(sinit)
	if err != nil {
		return nil, err
	}
	return sock, nil
}
