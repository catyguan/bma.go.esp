package espsocket

import (
	"bmautil/connutil"
	"bmautil/valutil"
	"esp/espnet/esnp"
	"net"
)

// ConnSocket
type ConnSocket struct {
	conn    *connutil.ConnExt
	m       connutil.ConnExtManager
	maxsize int
}

func NewConnSocket(conn *connutil.ConnExt, maxsize int) *ConnSocket {
	this := new(ConnSocket)
	this.conn = conn
	this.m = conn.Manager
	this.maxsize = maxsize
	return this
}

func NewConnSocketN(conn net.Conn, maxsize int) *ConnSocket {
	return NewConnSocket(connutil.NewConnExt(conn), maxsize)
}

func (this *ConnSocket) BaseConn() net.Conn {
	return this.conn
}

func (this *ConnSocket) String() string {
	return this.conn.String()
}

func (this *ConnSocket) GetProperty(name string) (interface{}, bool) {
	if name == PROP_MESSAGE_MAXSIZE {
		return this.maxsize, true
	}
	return this.conn.GetProperty(name)
}

func (this *ConnSocket) SetProperty(name string, val interface{}) bool {
	if name == PROP_MESSAGE_MAXSIZE {
		this.maxsize = valutil.ToInt(val, this.maxsize)
		return true
	}
	this.conn.SetProperty(name, val)
	return true
}

func (this *ConnSocket) AskFinish() {
	connutil.AskFinish(this.m, this.conn)
}

func (this *ConnSocket) AskClose() {
	connutil.AskClose(this.m, this.conn)
}

func (this *ConnSocket) WriteMessage(msg *esnp.Message) error {
	bs, err := msg.ToBytes()
	if err != nil {
		return err
	}
	_, err = this.conn.Write(bs)
	return err
}

func (this *ConnSocket) ReadMessage(decodeErr bool) (*esnp.Message, error) {
	reader := esnp.NewIODecodeReader(this.conn, nil)
	msg := esnp.NewMessage()
	err := msg.ReadAll(reader, this.maxsize)
	if err != nil {
		return nil, err
	}
	if decodeErr {
		err = msg.ToError()
	}
	return msg, err
}

func (this *ConnSocket) IsBreak() bool {
	ch := this.conn
	return ch.CheckBreak()
}
