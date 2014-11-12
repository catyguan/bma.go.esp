package auth

import (
	"acl"
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
)

const (
	SN_AUTH          = "espnet.auth"
	OP_AUTH          = "do"
	HEADER_AUTH_TYPE = "espnet.auth.type"
	HEADER_AUTH_FLAG = "espnet.auth.flag"
)

type NodeAuth func(sock *espsocket.Socket, msg *esnp.Message) (bool, *acl.User, error)
