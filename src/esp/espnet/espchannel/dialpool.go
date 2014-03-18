package espchannel

import (
	"bmautil/socket"
	"time"
)

// ChannelFactory
type dialPoolChannelFactory struct {
	service      *socket.DialPool
	channelCoder string
	getTimeout   time.Duration
}

func NewDialPoolChannelFactory(pool *socket.DialPool, chcoder string, getTimeout time.Duration) ChannelFactory {
	r := new(dialPoolChannelFactory)
	r.service = pool
	r.channelCoder = chcoder
	r.getTimeout = getTimeout
	return r
}

func (this *dialPoolChannelFactory) String() string {
	return this.service.String()
}

func (this *dialPoolChannelFactory) Start() bool {
	return this.service.Start()
}

func (this *dialPoolChannelFactory) Run() bool {
	return this.service.Run()
}

func (this *dialPoolChannelFactory) Close() bool {
	return this.service.Close()
}

func (this *dialPoolChannelFactory) NewChannel() (Channel, error) {
	sock, err := this.service.GetSocket(this.getTimeout, true)
	if err != nil {
		return nil, err
	}
	r := NewSocketChannel(sock, this.channelCoder)
	r.socketReturn = this.socketReturn
	return r, nil
}

func (this *dialPoolChannelFactory) socketReturn(s *socket.Socket) {
	this.service.ReturnSocket(s)
}

func (this *dialPoolChannelFactory) IsBreak() bool {
	return this.service.IsBreak()
}
