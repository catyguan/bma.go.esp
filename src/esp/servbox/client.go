package servbox

import (
	"bmautil/connutil"
	"esp/espnet/esnp"
	"esp/espnet/espservice"
	"esp/espnet/espsocket"
	"logger"
	"net"
	"time"
)

type Client struct {
	name     string
	info     string
	entry    espservice.ServiceEntry
	config   *clientInfo
	listener net.Listener
	services []string
	boxc     chan interface{}

	SkipKill bool
}

func NewClient(n string, info string, entry espservice.ServiceEntry) *Client {
	r := new(Client)
	r.name = n
	r.info = info
	r.entry = entry
	return r
}

func (this *Client) Add(serviceName string) {
	this.services = append(this.services, serviceName)
}

func (this *Client) accept(conn net.Conn) {
	sock := espsocket.NewConnSocketN(conn, espsocket.DEFAULT_MESSAGE_MAXSIZE)
	defer sock.AskClose()
	this.entry(sock)
}

func (this *Client) joinBox(cfg *clientInfo, addr string) {
	if this.boxc != nil {
		close(this.boxc)
	}
	ch := make(chan interface{}, 1)
	this.boxc = ch
	go func() {
		var lastErr string
		var lastTime time.Time
		var log bool
		n := 0
		for {
			err := this.doJoinBox(ch, cfg, addr)
			if err != nil {
				n++
				if err.Error() != lastErr || time.Since(lastTime) > 60*time.Second {
					lastErr = err.Error()
					lastTime = time.Now()
					log = true
					logger.Warn(tag, "join box fail - %d, %s", n, err)
				}
			}
			tm := time.NewTimer(1 * time.Second)
			select {
			case <-tm.C:
				// retry
				if log {
					log = false
					logger.Info(tag, "retry join box ...")
				}
			case <-ch:
				// close
				return
			}
		}
	}()
}

func (this *Client) doJoinBox(ch chan interface{}, cfg *clientInfo, addr string) error {
	conn, err0 := net.DialTimeout(cfg.BoxNet, cfg.BoxAddress, 5*time.Second)
	if err0 != nil {
		return err0
	}
	defer conn.Close()

	ce := connutil.NewConnExt(conn)
	ce.SetDeadline(time.Now().Add(5 * time.Second))
	sock := espsocket.NewConnSocketN(ce, 0)
	if true {
		var q objJoinQ
		q.Net = cfg.Net
		q.Address = addr
		q.NodeName = cfg.NodeName
		q.SerivceNames = this.services
		q.Info = this.info
		q.SkipKill = this.SkipKill
		msg := esnp.NewRequestMessage()
		esnp.MessageLineCoders.Address.Set(msg, esnp.ADDRESS_OP, op_Join)
		q.Encode(msg)
		err1 := sock.WriteMessage(msg)
		if err1 != nil {
			return err1
		}
		_, err2 := sock.ReadMessage(true)
		if err2 != nil {
			return err2
		}
	}
	if true {
		msg := esnp.NewRequestMessage()
		esnp.MessageLineCoders.Address.Set(msg, esnp.ADDRESS_OP, op_Active)
		err1 := sock.WriteMessage(msg)
		if err1 != nil {
			return err1
		}
		_, err2 := sock.ReadMessage(true)
		if err2 != nil {
			return err2
		}
	}
	ce.ClearDeadline()
	logger.Info(tag, "servBox(%s, %s) joined", cfg.BoxNet, cfg.BoxAddress)
	logger.Info(tag, "active services %v", this.services)

	for {
		if ce.CheckBreakDeadline(time.Now().Add(1 * time.Second)) {
			logger.Info(tag, "servBox(%s, %s) break", cfg.BoxNet, cfg.BoxAddress)
			return nil
		}
		select {
		case <-ch:
			// stop
			logger.Debug(tag, "joinBox stop")
			return nil
		default:
		}
	}

	return nil
}
