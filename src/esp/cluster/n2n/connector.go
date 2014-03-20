package n2n

import (
	"bmautil/socket"
	"esp/espnet/esnp"
	"esp/espnet/espchannel"
	"fmt"
	"logger"
)

type connector struct {
	service *Service
	name    string
	url     *esnp.URL

	pool *socket.DialPool
	ch   espchannel.Channel
}

func (this *connector) InitConnector(s *Service, n string, url *esnp.URL) {
	this.service = s
	this.name = n
	this.url = url
	go this.start()
}

func (this *connector) Close() {
	if this.pool != nil {
		this.pool.AskClose()
	}
	if this.ch != nil {
		this.ch.AskClose()
	}
}

func (this *connector) start() {
	host := this.url.GetHost()
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
	this.ch = espchannel.NewSocketChannel(sock, espchannel.SOCKET_CHANNEL_CODER_ESPNET)
	logger.Debug(tag, "%s connected", this.ch)
	this.service.goo.DoNow(func() {
		this.service.doChannelAccept(this.name, this.url, this.ch)
	})
	return nil
}
