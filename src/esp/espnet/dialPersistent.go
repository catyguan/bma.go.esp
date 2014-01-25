package espnet

import (
	"bmautil/socket"
	"time"
)

type DialPersistentAcceptor func(sock *socket.Socket) error

type DialPersistentConfig struct {
	Address   string
	TimeoutMS int
}

func NewDialPersistent(name string, cfg *DialPersistentConfig, acceptor DialPersistentAcceptor) (*DialPool, error) {
	dcfg := new(DialPoolConfig)
	dcfg.Dial.Address = cfg.Address
	dcfg.Dial.TimeoutMS = cfg.TimeoutMS
	dcfg.InitSize = 1
	dcfg.MaxSize = 1

	err := dcfg.Valid()
	if err != nil {
		return nil, err
	}
	var dial *DialPool
	dial = NewDialPool(name, dcfg, func(s *socket.Socket) error {
		go func() {
			sock, err := dial.GetSocket(time.Duration(cfg.TimeoutMS)*time.Millisecond, true)
			if err != nil {
				return
			}
			err = acceptor(sock)
			if err != nil {
				sock.Close()
			}
		}()
		return nil
	})
	return dial, nil
}
