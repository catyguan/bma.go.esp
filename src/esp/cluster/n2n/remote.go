package n2n

import (
	"esp/cluster/nodebase"
	"esp/espnet/espsocket"
	"fmt"
)

type remoteInfo struct {
	service *Service

	nodeId   nodebase.NodeId
	nodeName string
	nodeHost string

	sock  *espsocket.Socket
	token string
}

func (this *remoteInfo) String() string {
	return fmt.Sprintf("%d(%s,%s)", this.nodeId, this.nodeName, this.nodeHost)
}

func (this *remoteInfo) InitRemoteInfo(s *Service, id nodebase.NodeId, name string, host string, sock *espsocket.Socket) error {
	this.service = s
	this.nodeId = id
	this.nodeName = name
	this.nodeHost = host

	sock.SetCloseListener("n2n.service", func() {
		this.service.goo.DoNow(func() {
			this.service.doRemoteClosed(this)
		})
	})
	this.sock = sock
	return nil
}

func (this *remoteInfo) close(removeCloseListener bool) {
	if removeCloseListener {
		this.sock.SetCloseListener("n2n.service", nil)
	}
	this.sock.AskClose()
}
