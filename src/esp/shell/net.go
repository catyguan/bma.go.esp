package shell

import (
	"bmautil/netutil"
)

type NetWriter struct {
	channel *netutil.Channel
}

func NewNetWriter(ch *netutil.Channel) *NetWriter {
	r := new(NetWriter)
	r.channel = ch
	return r
}

func (this *NetWriter) Write(msg string) bool {
	this.channel.Write([]byte(msg))
	return true
}

func (this *NetWriter) Close() {
	this.channel.CloseChannel()
}

func (this *NetWriter) AddCloseListner(lis func()) {
	this.channel.AddListener(lis)
}
