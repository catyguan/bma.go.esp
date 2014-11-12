package n2n

import (
	"acl"
	"crypto/md5"
	"esp/cluster/nodebase"
	"esp/espnet/auth"
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"fmt"
	"logger"
	"math/rand"
	"time"
)

type authReq struct {
	NodeId nodebase.NodeId
	Token  string
}

func (this *authReq) String() string {
	return fmt.Sprintf("[NodeId=%d, Token=%s]", this.NodeId, this.Token)
}

func (this *authReq) Write(msg *esnp.Message) error {
	xd := msg.XDatas()
	xd.Add(1, this.NodeId, nodebase.NodeIdCoder)
	xd.Add(2, this.Token, esnp.Coders.String)
	return nil
}

func (this *authReq) Read(msg *esnp.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		switch it.Xid() {
		case 1:
			v, err := it.Value(nodebase.NodeIdCoder)
			if err != nil {
				return err
			}
			if v != nil {
				this.NodeId = v.(nodebase.NodeId)
			}
		case 2:
			v, err := it.Value(esnp.Coders.String)
			if err != nil {
				return err
			}
			if v != nil {
				this.Token = v.(string)
			}
		}
	}
	return nil
}

type authToken struct {
	Token string
}

func (this *authToken) String() string {
	return fmt.Sprintf("[Token=%s]", this.Token)
}

func (this *authToken) Write(msg *esnp.Message) error {
	xd := msg.XDatas()
	xd.Add(1, this.Token, esnp.Coders.String)
	return nil
}

func (this *authToken) Read(msg *esnp.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		switch it.Xid() {
		case 1:
			v, err := it.Value(esnp.Coders.String)
			if err != nil {
				return err
			}
			if v != nil {
				this.Token = v.(string)
			}
		}
	}
	return nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randSeq(n int) string {
	l := len(letters)
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(l)]
	}
	return string(b)
}

func code(code, token string) string {
	h := md5.New()
	h.Write([]byte(code + ":" + token))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (this *Service) NodeAuth(sock *espsocket.Socket, msg *esnp.Message) (bool, *acl.User, error) {
	at := auth.GetAuthType(msg)
	if at != "n2n" {
		return false, nil, nil
	}

	req := new(authReq)
	err0 := req.Read(msg)
	if err0 != nil {
		return false, nil, err0
	}

	token := req.Token
	socktk, ok := sock.GetProperty("auth.code")
	if !ok {
		token = ""
	}
	if token == "" {
		tk := randSeq(16)
		sock.SetProperty("auth.code", tk)
		rmsg := msg.ReplyMessage()
		atk := new(authToken)
		atk.Token = tk
		atk.Write(rmsg)
		sock.PostMessage(rmsg)
		return true, nil, nil
	}

	mytoken := code(this.config.Code, socktk.(string))
	sock.SetProperty("auth.code", nil)
	if mytoken != token {
		logger.Warn(tag, "auth token fail %s != %s", mytoken, token)
		return false, nil, fmt.Errorf("auth token fail")
	}

	host, _ := sock.GetRemoteAddr()
	user := acl.NewUser(fmt.Sprintf("%d", req.NodeId), host, nil)
	logger.Debug(tag, "'%s' auth -> %s", sock, user)
	return true, user, nil
}

func (this *Service) PostAuth(sock *espsocket.Socket, cd string) error {
	token := ""
	if true {
		msg := esnp.NewMessage()
		addr := msg.GetAddress()
		addr.SetService(auth.SN_AUTH)
		addr.SetOp(auth.OP_AUTH)
		auth.SetAuthType(msg, "n2n")

		req := new(authReq)
		req.NodeId = nodebase.Id
		req.Token = ""
		req.Write(msg)

		rmsg, err := sock.Call(msg, 1*time.Second)
		if err != nil {
			return err
		}

		atk := new(authToken)
		atk.Read(rmsg)
		token = code(cd, atk.Token)
		// fmt.Println("FUCK @@@@@@@@@", cd, atk.Token, token)
	}

	if true {
		msg := esnp.NewMessage()
		addr := msg.GetAddress()
		addr.SetService(auth.SN_AUTH)
		addr.SetOp(auth.OP_AUTH)
		auth.SetAuthType(msg, "n2n")

		req := new(authReq)
		req.NodeId = nodebase.Id
		req.Token = token
		req.Write(msg)

		rmsg, err := sock.Call(msg, 1*time.Second)
		if err != nil {
			return err
		}
		_, ok := auth.CheckAuth(rmsg)
		if !ok {
			return fmt.Errorf("auth return fail")
		}
	}

	return nil

}
