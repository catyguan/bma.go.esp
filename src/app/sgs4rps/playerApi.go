package sgs4rps

import (
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
)

// api
type playerAPI interface {
	NewPlayer(sock *espsocket.Socket) (int, error)
	SetNick(psid int, n string) error
	Play(psid int, rps int) error
	Quit(psid int) error
}

// message packer
const (
	op_SetNick = "SetNick"
	op_Play    = "Play"
)

type packPlayerAPI struct {
	serviceName string
}

func (this *packPlayerAPI) InitService(sname string) {
	this.serviceName = sname
}

func (this *packPlayerAPI) packSetNick(n string) (*esnp.Message, error) {
	r := esnp.NewMessage()
	r.SetFlag(esnp.FLAG_INFO)
	r.GetAddress().SetCall(this.serviceName, op_SetNick)
	r.Datas().Set("name", n)
	return r, nil
}

func (this *packPlayerAPI) unpackSetNick(m *esnp.Message) (string, error) {
	n, err1 := m.Datas().GetString("name", "")
	if err1 != nil {
		return "", err1
	}
	return n, nil
}

func (this *packPlayerAPI) packPlay(rps int) (*esnp.Message, error) {
	r := esnp.NewMessage()
	r.SetFlag(esnp.FLAG_INFO)
	r.GetAddress().SetCall(this.serviceName, op_SetNick)
	r.Datas().Set("rps", rps)
	return r, nil
}

func (this *packPlayerAPI) unpackPlay(m *esnp.Message) (int, error) {
	rps, err1 := m.Datas().GetInt("rps", 0)
	if err1 != nil {
		return "", err1
	}
	return rps, nil
}
