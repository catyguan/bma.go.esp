package n2n

import (
	"esp/cluster/nodeinfo"
	"esp/espnet/esnp"
	"esp/espnet/espchannel"
	"esp/espnet/espterminal"
	"esp/espnet/esptunnel"
	"fmt"
)

type remoteInfo struct {
	service *Service

	nodeId   nodeinfo.NodeId
	nodeName string
	nodeURL  *esnp.URL

	tunnel *esptunnel.Tunnel
	tm     *espterminal.Terminal
}

func (this *remoteInfo) String() string {
	return fmt.Sprintf("%d(%s)", this.nodeId, this.nodeName)
}

func (this *remoteInfo) InitRemoteInfo(s *Service, id nodeinfo.NodeId, name string, url string) error {
	v, err := esnp.ParseURL(url)
	if err != nil {
		return err
	}

	this.service = s
	this.nodeId = id
	this.nodeName = name
	this.nodeURL = v

	this.tunnel = esptunnel.NewTunnel(this.nodeName)
	this.tunnel.CloseOnBreak = true
	this.tunnel.SetCloseListener("this", func() {
		this.service.goo.DoNow(func() {
			this.service.doRemoteClosed(this)
		})
	})

	this.tm = new(espterminal.Terminal)
	this.tm.InitTerminal(name)
	this.tm.SetMessageListner(func(msg *esnp.Message) error {
		return this.service.Serve(this.tunnel, msg)
	})
	this.tm.Connect(this.tunnel)
	return nil
}

func (this *remoteInfo) Add(ch espchannel.Channel) {
	this.tunnel.Add(ch)
}

func (this *remoteInfo) Close() {
	if this.tunnel != nil {
		this.tunnel.AskClose()
	}
}
