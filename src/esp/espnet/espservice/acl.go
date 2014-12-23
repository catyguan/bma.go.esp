package espservice

import (
	"acl"
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"net"
)

const (
	PROP_ACL_USER = "acl.user"
)

func GetUser(sock espsocket.Socket) *acl.User {
	v, ok := sock.GetProperty(PROP_ACL_USER)
	if ok {
		if r, ok2 := v.(*acl.User); ok2 {
			return r
		}
	}
	return nil
}

func SetUser(sock espsocket.Socket, user *acl.User) bool {
	return sock.SetProperty(PROP_ACL_USER, user)
}

type AclServiceHandler struct {
	name string
	h    ServiceHandler
}

func NewAclServiceHandler(n string, h ServiceHandler) *AclServiceHandler {
	r := new(AclServiceHandler)
	r.name = n
	r.h = h
	return r
}

func (this *AclServiceHandler) Serve(sock espsocket.Socket, msg *esnp.Message) error {
	var user *acl.User
	if tmp, ok := sock.GetProperty(PROP_ACL_USER); ok {
		user, ok = tmp.(*acl.User)
	}
	if user == nil {
		if v, ok := espsocket.GetProperty(sock, espsocket.PROP_SOCKET_REMOTE_ADDR); ok {
			if str, ok := v.(string); ok {
				ip, _, _ := net.SplitHostPort(str)
				user = acl.NewUser("anonymous", ip, nil)
				sock.SetProperty(PROP_ACL_USER, user)
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
			return err
		}
	}
	return this.h(sock, msg)
}
