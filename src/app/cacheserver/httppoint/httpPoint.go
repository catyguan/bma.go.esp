package httppoint

import (
	"app/cacheserver"
	"bmautil/valutil"
	"boot"
	"config"
	"encoding/json"
	"httpserver"
	"logger"
	"net/http"
	"strings"
)

const (
	tag = "httpPoint"
)

type HttpPoint struct {
	name    string
	service *cacheserver.CacheService

	configInfo *configInfo
	server     *httpserver.HttpServer
	handler    *http.ServeMux
}

func NewHttpPoint(name string, s *cacheserver.CacheService) *HttpPoint {
	this := new(HttpPoint)
	this.name = name
	this.service = s
	this.server = httpserver.NewHttpServer(name, this)
	return this
}

func (this *HttpPoint) Name() string {
	return this.name
}

type configInfo struct {
	httpserver.HttpServerConfigInfo
	AppPath string
	Disable bool
}

func (this *HttpPoint) Init() bool {
	cfg := new(configInfo)
	if config.GetBeanConfig(this.name, cfg) {
		if !this.server.InitConfig(&cfg.HttpServerConfigInfo) {
			return false
		}
		if cfg.AppPath == "" {
			cfg.AppPath = "/"
		}
		if !strings.HasSuffix(cfg.AppPath, "/") {
			cfg.AppPath += "/"
		}
	} else {
		cfg.Disable = true
	}
	this.configInfo = cfg

	if this.configInfo.Disable {
		logger.Debug(tag, "disable")
		return true
	}
	this.initServeMux()
	return true
}

func (this *HttpPoint) Start() bool {
	if this.configInfo.Disable {
		return true
	}
	return this.server.Start()
}

func (this *HttpPoint) Run() bool {
	if this.configInfo.Disable {
		return true
	}
	return this.server.Run()
}

func (this *HttpPoint) Stop() bool {
	if this.configInfo.Disable {
		return true
	}
	return this.server.Stop()
}

func (this *HttpPoint) DefaultBoot() {
	boot.Define(boot.INIT, this.name, this.Init)
	boot.Define(boot.START, this.name, this.Start)
	boot.Define(boot.RUN, this.name, this.Run)
	boot.Define(boot.STOP, this.name, this.Stop)

	boot.Install(this.name, this)
}

func (this *HttpPoint) initServeMux() {
	m := http.NewServeMux()
	path := this.configInfo.AppPath
	m.HandleFunc(path+"get", this.serveGet)
	m.HandleFunc(path+"load", this.serveGet)
	m.HandleFunc(path+"put", this.serveGet)
	m.HandleFunc(path+"erase", this.serveErase)
	this.handler = m
}

func (this *HttpPoint) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger.Debug(tag, "serveHTTP %s", req.RequestURI)
	this.handler.ServeHTTP(w, req)
}

type resultInfo struct {
	doneInfo
	Value  string
	Traces []string
}

type doneInfo struct {
	Done bool
	Err  string
}

func (this *HttpPoint) serveGet(w http.ResponseWriter, r *http.Request) {

	groupName := r.FormValue("g")
	key := r.FormValue("k")
	timeout := valutil.ToInt32(r.FormValue("t"), 0)
	trace := valutil.ToBool(r.FormValue("d"), false)
	notLoad := valutil.ToBool(r.FormValue("n"), false)

	logger.Debug(tag, "serveGet(%v,%v)", groupName, key)

	req := cacheserver.NewGetRequest(key)
	if timeout == 0 {
		req.TimeoutMs = 5 * 1000
	} else {
		req.TimeoutMs = timeout * 1000
	}
	req.NotLoad = notLoad
	req.Trace = trace

	rep := make(chan *cacheserver.GetResult, 1)
	defer close(rep)

	err := this.service.Get(groupName, req, rep)
	if err != nil {
		logger.Warn(tag, "CacheServerGet fail - %s", err.Error())
		this.writeDone(w, err)
		return
	}

	result := <-rep
	if result == nil {
		err = logger.Warn(tag, "CacheServerGet null result return")
		this.writeDone(w, err)
		return
	}

	logger.Debug(tag, "CacheServerGet %v -> %v", req, result.Done)

	tr := resultInfo{}
	tr.Done = result.Done
	if result.Value != nil {
		tr.Value = string(result.Value)
	}
	if result.Err != nil {
		tr.Err = result.Err.Error()
	}
	tr.Traces = result.TraceInfo

	out, err := json.Marshal(tr)
	if err != nil {
		msg := logger.Warn(tag, "format result fail %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg.Error()))
	} else {
		w.Write([]byte(out))
	}
}

func (this *HttpPoint) writeDone(w http.ResponseWriter, err error) {
	tr := doneInfo{}
	if err != nil {
		tr.Done = false
		tr.Err = err.Error()
	} else {
		tr.Done = true
	}
	out, err := json.Marshal(tr)
	if err != nil {
		msg := logger.Warn(tag, "format result fail %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg.Error()))
	} else {
		w.Write([]byte(out))
	}
}

func (this *HttpPoint) serveLoad(w http.ResponseWriter, r *http.Request) {
	groupName := r.FormValue("g")
	key := r.FormValue("k")
	logger.Debug(tag, "serveLoad(%s,%s)", groupName, key)
	err := this.service.Load(groupName, key)
	if err != nil {
		logger.Warn(tag, "serveLoad fail - %s", err.Error())
	}
	this.writeDone(w, err)
}

func (this *HttpPoint) servePut(w http.ResponseWriter, r *http.Request) {
	groupName := r.FormValue("g")
	key := r.FormValue("k")
	value := r.FormValue("v")
	logger.Debug(tag, "servePut(%s,%s,%s)", groupName, key, value)
	err := this.service.Put(groupName, key, []byte(value), 0)
	if err != nil {
		logger.Warn(tag, "servePut fail - %s", err.Error())
	}
	this.writeDone(w, err)
}

func (this *HttpPoint) serveErase(w http.ResponseWriter, r *http.Request) {
	groupName := r.FormValue("g")
	key := r.FormValue("k")

	logger.Debug(tag, "serveErase(%s,%s)", groupName, key)

	_, err := this.service.Delete(groupName, key)
	if err != nil {
		logger.Warn(tag, "serveErase fail - %s", err.Error())
	}
	this.writeDone(w, err)
}
