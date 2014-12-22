package aclesnpmux

import (
	"acl"
	"esp/espnet/esnp"
	"esp/espnet/espservice"
	"esp/espnet/espsocket"
	"net"
)

func GetUser(sock espsocket.Socket) *acl.User {
	v, ok := sock.GetProperty("acl.user")
	if ok {
		if r, ok2 := v.(*acl.User); ok2 {
			return r
		}
	}
	return nil
}

func SetUser(sock espsocket.Socket, user *acl.User) bool {
	return sock.SetProperty("acl.user", user)
}

type AclServerMux struct {
	name string
	h    espservice.ServiceHandler
}

func NewAclServerMux(n string, h espservice.ServiceHandler) *AclServerMux {
	r := new(AclServerMux)
	r.name = n
	r.h = h
	return r
}

func (this *AclServerMux) DoServe(sock espsocket.Socket, msg *esnp.Message) error {
	var user *acl.User
	if tmp, ok := sock.GetProperty("acl.user"); ok {
		user, ok = tmp.(*acl.User)
	}
	if user == nil {
		if v, ok := espsocket.GetProperty(sock, espsocket.PROP_SOCKET_REMOTE_ADDR); ok {
			if str, ok := v.(string); ok {
				ip, _, _ := net.SplitHostPort(str)
				user = acl.NewUser("anonymous", ip, nil)
				sock.SetProperty("acl.user", user)
			}
		}
	}

	addr := msg.GetAddress()
	if addr != nil {
		sname := addr.GetService()
		opname := addr.GetOp()
		var ps []string
		ps = append(ps, this.name)
		if sname != "" {
			ps = append(ps, sname)
		} else {
			ps = append(ps, "unknow")
		}
		if opname != "" {
			ps = append(ps, opname)
		}
		if user == nil {
			user = acl.NewUser("anonymous", "unknow", nil)
		}
		err := acl.Assert(user, ps, nil)
		if err != nil {
			rmsg := msg.ReplyMessage()
			rmsg.BeError(err)
			sock.WriteMessage(rmsg)
			return err
		}
	}
	return this.h(sock, msg)
}
