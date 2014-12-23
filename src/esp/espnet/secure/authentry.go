package secure

import (
	"esp/espnet/espservice"
	"esp/espnet/espsocket"
	"time"
)

type BaseSecureConfig struct {
	MaxAuthTime    time.Duration
	NotAuthMaxSize int
	AuthMaxSize    int
	Entry          espservice.ServiceEntry
}

func (this *BaseSecureConfig) InitDefault() {
	this.NotAuthMaxSize = 4 * 1024
	this.AuthMaxSize = espsocket.DEFAULT_MESSAGE_MAXSIZE
	this.MaxAuthTime = 5 * time.Second
}

func (this *BaseSecureConfig) Begin(sock espsocket.Socket) {
	espsocket.SetDeadline(sock, time.Now().Add(this.MaxAuthTime))
	sock.SetProperty(espsocket.PROP_MESSAGE_MAXSIZE, this.NotAuthMaxSize)
}

func (this *BaseSecureConfig) DoNext(sock espsocket.Socket) {
	espsocket.ClearDeadline(sock)
	sock.SetProperty(espsocket.PROP_MESSAGE_MAXSIZE, this.AuthMaxSize)
	if this.Entry != nil {
		this.Entry(sock)
	}
}
