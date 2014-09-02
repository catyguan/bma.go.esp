package main

import (
	"bmautil/socket"
	"bmautil/valutil"
	"encoding/json"
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"fmt"
	"logger"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type SockInfo struct {
	sock  *espsocket.Socket
	timer *time.Timer
}

type Service struct {
	name   string
	config *configInfo
	robj   *rand.Rand

	slock   sync.RWMutex
	sockets map[string]*SockInfo
}

func (this *Service) closeSocket(sid string) {
	this.slock.Lock()
	defer this.slock.Unlock()
	si := this.sockets[sid]
	if si != nil {
		logger.Debug(tag, "socket close - %s", sid)
		delete(this.sockets, sid)
		si.sock.AskClose()
	}
}

func (this *Service) closeAll() {
	this.slock.Lock()
	defer this.slock.Unlock()
	for k, si := range this.sockets {
		delete(this.sockets, k)
		si.sock.AskClose()
	}
}

func (this *Service) InvokeESNP(w http.ResponseWriter, req *http.Request) {
	var sid string
	if cookie, err := req.Cookie("sid"); err == nil {
		sid = cookie.Value
	}
	if sid == "" {
		sid = fmt.Sprintf("%x", this.robj.Uint32())
		logger.Debug(tag, "new sid = %s", sid)

		cookie := new(http.Cookie)
		cookie.Name = "sid"
		cookie.Value = sid
		cookie.Path = this.config.CookiePath
		http.SetCookie(w, cookie)
	} else {
		logger.Debug(tag, "sid = %s", sid)
	}

	path := req.URL.Path
	ps := strings.SplitN(path, "/", 4)
	if len(ps) < 4 {
		http.Error(w, "/i/serviceName/opName", http.StatusBadRequest)
		return
	}

	sn := ps[2]
	op := ps[3]
	req.ParseForm()
	form := req.Form

	msg := esnp.NewRequestMessage()
	msg.GetAddress().SetCall(sn, op)
	ds := msg.Datas()
	for k, _ := range form {
		var rv interface{}
		fv := form.Get(k)
		rv = fv
		ks := strings.SplitN(k, "-", 2)
		if len(ks) == 2 {
			k = ks[0]
			t := strings.ToLower(ks[1])
			switch t {
			case "b", "bool":
				rv = valutil.ToBool(fv, false)
			case "i", "int", "int32":
				rv = valutil.ToInt(fv, 0)
			case "i8", "byte":
				rv = valutil.ToByte(fv, 0)
			case "i16", "short":
				rv = valutil.ToInt16(fv, 0)
			case "i64", "long":
				rv = valutil.ToInt64(fv, 0)
			case "f32", "float":
				rv = valutil.ToFloat32(fv, 0)
			case "f64", "double":
				rv = valutil.ToFloat64(fv, 0)
			case "a":
				a := make([]interface{}, 0)
				err := json.Unmarshal([]byte(fv), &a)
				if err != nil {
					rv = nil
					logger.Debug(tag, "parse array fail - %s", err)
				} else {
					rv = a
				}
			case "o":
				o := make(map[string]interface{}, 0)
				err := json.Unmarshal([]byte(fv), &o)
				if err != nil {
					rv = nil
					logger.Debug(tag, "parse array fail - %s", err)
				} else {
					rv = o
				}
			}
		}
		logger.Debug(tag, "data %s = %v", k, rv)
		if rv != nil {
			ds.Set(k, rv)
		}
	}

	this.slock.RLock()
	si := this.sockets[sid]
	this.slock.RUnlock()
	if si == nil {
		ok := func() bool {
			cfg := new(socket.DialConfig)
			cfg.Address = this.config.EsnpAddress
			sock, err := espsocket.Dial("sock"+sid, cfg, "")
			if err != nil {
				http.Error(w, fmt.Sprintf("connect %s fail - %s", this.config.EsnpAddress, err), http.StatusInternalServerError)
				return false
			}
			this.slock.Lock()
			defer this.slock.Unlock()
			osi := this.sockets[sid]
			if osi != nil {
				sock.AskClose()
			} else {
				si = new(SockInfo)
				si.sock = sock
				si.sock.SetCloseListener("", func() {
					this.closeSocket(sid)
				})
				si.timer = time.AfterFunc(time.Duration(this.config.ExpiresSec)*time.Second, func() {
					this.closeSocket(sid)
				})
				this.sockets[sid] = si
			}
			return true
		}()
		if !ok {
			return
		}
	}

	si.timer.Reset(time.Duration(this.config.ExpiresSec) * time.Second)
	sock := si.sock
	rmsg, err := sock.Call(msg, time.Duration(this.config.TimeoutMS)*time.Millisecond)
	if err != nil {
		http.Error(w, fmt.Sprintf("call %s::%s fail - %s", sn, op, err), http.StatusInternalServerError)
		return
	}
	ds2 := rmsg.Datas()
	ns := ds2.List()
	r := make(map[string]interface{})
	for _, n := range ns {
		var err1 error
		r[n], err1 = ds2.Get(n)
		if err1 != nil {
			http.Error(w, fmt.Sprintf("format result fail - %s", err1), http.StatusInternalServerError)
			return
		}
	}
	content, err2 := json.Marshal(r)
	if err2 != nil {
		http.Error(w, fmt.Sprintf("format result fail - %s", err2), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(content))
}
