package n2n

import (
	"esp/cluster/nodeinfo"
	"esp/espnet/esnp"
	"esp/espnet/espchannel"
	"esp/espnet/esptunnel"
)

type remoteInfo struct {
	service *Service

	nodeId   nodeinfo.NodeId
	nodeName string
	nodeURL  *esnp.URL

	tunnel *esptunnel.Tunnel
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
	return nil
}

func (this *remoteInfo) Add(ch espchannel.Channel) {
	if this.tunnel == nil {
		this.tunnel = esptunnel.NewTunnel(this.nodeName)
	}
	this.tunnel.Add(ch)
	ch.SetMessageListner(func(msg *esnp.Message) error {
		return this.service.Serve(ch, msg)
	})
}

func (this *remoteInfo) Close() {
	if this.tunnel != nil {
		this.tunnel.AskClose()
	}
}
