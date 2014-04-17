package main

import (
	"bmautil/socket"
	"bmautil/valutil"
	"boot"
	"encoding/json"
	"esp/espnet/esnp"
	"esp/espnet/espchannel"
	"esp/espnet/espclient"
	"fmt"
	"logger"
	"net/http"
	"strings"
	"time"
)

type configInfo struct {
	EsnpAddress string
	TimeoutMS   int
}

func (this *configInfo) Valid() error {
	if this.EsnpAddress == "" {
		return fmt.Errorf("EsnpAddress invalid")
	}
	if this.TimeoutMS <= 0 {
		this.TimeoutMS = 5000
	}
	return nil
}

func (this *configInfo) Compare(old *configInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	if this.EsnpAddress != old.EsnpAddress {
		return boot.CCR_NEED_START
	}
	if this.TimeoutMS != old.TimeoutMS {
		return boot.CCR_NEED_START
	}
	return boot.CCR_NONE
}

type Service struct {
	name   string
	config *configInfo

	client *espclient.ChannelClient
}

func (this *Service) Name() string {
	return this.name
}

func (this *Service) Prepare() {
}
func (this *Service) CheckConfig(ctx *boot.BootContext) bool {
	co := ctx.Config
	cfg := new(configInfo)
	if !co.GetBeanConfig(this.name, cfg) {
		logger.Error(tag, "'%s' miss config", this.name)
		return false
	}
	if err := cfg.Valid(); err != nil {
		logger.Error(tag, "'%s' config error - %s", this.name, err)
		return false
	}
	ccr := boot.NewConfigCheckResult(cfg.Compare(this.config), cfg)
	ctx.CheckFlag = ccr
	return true
}

func (this *Service) Init(ctx *boot.BootContext) bool {
	ccr := ctx.CheckResult()
	if ccr.Type == boot.CCR_NONE {
		return true
	}
	this.config = ccr.Config.(*configInfo)
	return true
}

func (this *Service) Start(ctx *boot.BootContext) bool {
	return true
}

func (this *Service) Run(ctx *boot.BootContext) bool {
	return true
}

func (this *Service) GraceStop(ctx *boot.BootContext) bool {
	return true
}

func (this *Service) Stop() bool {
	this.closeClient()
	return true
}

func (this *Service) Close() bool {
	return true
}

func (this *Service) Cleanup() bool {
	return true
}

func (this *Service) closeClient() {
	if this.client != nil {
		this.client.Close()
		this.client = nil
	}
}

func (this *Service) InvokeESNP(w http.ResponseWriter, req *http.Request) {
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
		rv = form.Get(k)
		ks := strings.SplitN(k, "-", 2)
		if len(ks) == 2 {
			k = ks[0]
			t := strings.ToLower(ks[1])
			v := rv
			switch t {
			case "b", "bool":
				rv = valutil.ToBool(v, false)
			case "i", "int", "int32":
				rv = valutil.ToInt(v, 0)
			case "i8", "byte":
				rv = valutil.ToByte(v, 0)
			case "i16", "short":
				rv = valutil.ToInt16(v, 0)
			case "i64", "long":
				rv = valutil.ToInt64(v, 0)
			case "f32", "float":
				rv = valutil.ToFloat32(v, 0)
			case "f64", "double":
				rv = valutil.ToFloat64(v, 0)
			}
		}
		logger.Debug(tag, "data %s = %v", k, rv)
		ds.Set(k, rv)
	}

	if this.client == nil {
		c := espclient.NewChannelClient()
		cfg := new(socket.DialConfig)
		cfg.Address = this.config.EsnpAddress
		err := c.Dial(tag, cfg, espchannel.SOCKET_CHANNEL_CODER_ESPNET)
		if err != nil {
			http.Error(w, fmt.Sprintf("connect %s fail - %s", this.config.EsnpAddress, err), http.StatusInternalServerError)
			return
		}
		this.client = c
	}
	c := this.client

	rmsg, err := c.Call(msg, time.NewTimer(time.Duration(this.config.TimeoutMS)*time.Millisecond))
	if err != nil {
		this.closeClient()
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
