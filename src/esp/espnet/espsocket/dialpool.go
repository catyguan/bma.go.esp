package espsocket

/*
// dialPoolSocketFactory
type dialPoolSocketFactory struct {
	service      *socket.DialPool
	channelCoder string
	getTimeout   time.Duration
}

func NewDialPoolSocketFactory(pool *socket.DialPool, chcoder string, getTimeout time.Duration) SocketFactory {
	r := new(dialPoolSocketFactory)
	r.service = pool
	r.channelCoder = chcoder
	r.getTimeout = getTimeout
	return r
}

func (this *dialPoolSocketFactory) String() string {
	return this.service.String()
}

func (this *dialPoolSocketFactory) Start() bool {
	return this.service.Start()
}

func (this *dialPoolSocketFactory) Run() bool {
	return this.service.Run()
}

func (this *dialPoolSocketFactory) Close() bool {
	return this.service.Close()
}

func (this *dialPoolSocketFactory) NewSocket() (*Socket, error) {
	sock, err := this.service.GetSocket(this.getTimeout, true)
	if err != nil {
		return nil, err
	}
	ch := NewSocketChannel(sock, this.channelCoder)
	ch.socketReturn = this.socketReturn
	return NewSocket(ch), nil
}

func (this *dialPoolSocketFactory) socketReturn(s *socket.Socket) {
	this.service.ReturnSocket(s)
}

func (this *dialPoolSocketFactory) IsBreak() bool {
	return this.service.IsBreak()
}
*/
