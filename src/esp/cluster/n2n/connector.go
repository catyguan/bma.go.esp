package n2n

import (
	"bmautil/socket"
	"esp/espnet/espsocket"
	"fmt"
	"logger"
)

type connector struct {
	service *Service
	name    string
	host    string
	code    string

	pool *socket.DialPool
	sock *espsocket.Socket
}

func (this *connector) InitConnector(s *Service, n string, host string, code string) {
	this.service = s
	this.name = n
	this.host = host
	this.code = code
	go this.start()
}

func (this *connector) Close() {
	if this.pool != nil {
		this.pool.AskClose()
	}
	if this.sock != nil {
		this.sock.AskClose()
	}
}

func (this *connector) start() {
	host := this.host
	logger.Debug(tag, "connector[%s] start -> %s", this.name, host)
	cfg := new(socket.DialPoolConfig)
	cfg.Dial.Address = host
	cfg.MaxSize = 1
	cfg.InitSize = 1
	this.pool = socket.NewDialPool(fmt.Sprintf("remote-%s-pool", this.name), cfg, nil)
	this.pool.Start()
	this.pool.Run()
	ss := socket.NewSocketServer4Dial(this.pool, 10)
	ss.SetAcceptor(this.onSocketAccept)
}

func (this *connector) onSocketAccept(sock *socket.Socket) error {
	ch := espsocket.NewSocketChannel(sock, espsocket.SOCKET_CHANNEL_CODER_ESPNET)
	this.sock = espsocket.NewSocket(ch)
	logger.Debug(tag, "%s connected", this.sock)
	this.service.goo.DoNow(func() {
		this.service.doSocketAccept(this.name, this.host, this.code, this.sock)
	})
	return nil
}
