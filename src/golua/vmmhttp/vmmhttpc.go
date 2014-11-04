package vmmhttp

import (
	"bmautil/httputil"
	"bmautil/valutil"
	"esp/acclog"
	"fmt"
	"golua"
	"io"
	"io/ioutil"
	"logger"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func InitGoLuaWithHttpClient(gl *golua.GoLua, accLog *acclog.Service, accName string) {
	HttpClientModule(accLog, accName).Bind(gl)
}

type httpClientModule struct {
	accLog  *acclog.Service
	accName string
}

func HttpClientModule(accLog *acclog.Service, accName string) *golua.VMModule {
	m := golua.NewVMModule("httpclient")
	mo := &httpClientModule{accLog, accName}
	m.Init("exec", &GOF_httpclient_exec{mo})
	m.Init("getContent", &GOF_httpclient_getContent{mo})
	return m
}

type httpclientRequest struct {
	URL       string
	Headers   map[string]string
	Post      bool
	Data      map[string]interface{}
	TimeoutMS int
}

func (this *httpclientRequest) Valid() error {
	if this.URL == "" {
		return fmt.Errorf("url empty")
	}
	if this.TimeoutMS <= 0 {
		this.TimeoutMS = 5 * 1000
	}
	return nil
}

func (this *httpClientModule) doExecute(vm *golua.VM, req *httpclientRequest) (map[string]interface{}, error) {
	err0 := req.Valid()
	if err0 != nil {
		return nil, err0
	}

	var body io.Reader
	method := "GET"
	qurl := req.URL
	data := make(url.Values)
	for k, v := range req.Data {
		data.Add(k, valutil.ToString(v, ""))
	}
	if req.Post {
		method = "POST"
		body = strings.NewReader(data.Encode())
	} else {
		if strings.Contains(qurl, "?") {
			qurl = qurl + "&" + data.Encode()
		} else {
			qurl = qurl + "?" + data.Encode()
		}
	}
	hreq, err2 := http.NewRequest(method, qurl, body)
	if err2 != nil {
		return nil, err2
	}
	if req.Post {
		hreq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	for k, v := range req.Headers {
		hreq.Header.Set(k, v)
	}
	client := httputil.NewHttpClient(time.Millisecond * time.Duration(req.TimeoutMS))
	logger.Debug(tag, "[%s] http '%s' start", vm, qurl)

	ts := time.Now()
	hresp, err3 := client.Do(hreq)
	te := time.Now()
	if err3 != nil {
		logger.Debug(tag, "[%s] http '%s' fail '%s'", vm, qurl, err3)
		if this.accLog != nil {
			ainfo := make(map[string]interface{})
			ainfo["url"] = qurl
			ainfo["error"] = err3.Error()
			ar := acclog.NewCommonLog(ainfo, te.Sub(ts).Seconds())
			this.accLog.Write(this.accName, ar)
		}
		return nil, err3
	}
	logger.Debug(tag, "[%s] http '%s' end '%d'", vm, qurl, hresp.StatusCode)
	defer hresp.Body.Close()
	respBody, err4 := ioutil.ReadAll(hresp.Body)
	if err4 != nil {
		if this.accLog != nil {
			ainfo := make(map[string]interface{})
			ainfo["url"] = qurl
			ainfo["error"] = err4.Error()
			ar := acclog.NewCommonLog(ainfo, te.Sub(ts).Seconds())
			this.accLog.Write(this.accName, ar)
		}
		return nil, err4
	}
	m := make(map[string]interface{})
	m["Status"] = hresp.StatusCode
	hs := make(map[string]string)
	for k, _ := range hresp.Header {
		v := hresp.Header.Get(k)
		hs[k] = v
	}
	m["Header"] = hs
	content := string(respBody)
	m["Content"] = content
	m["Time"] = te.Sub(ts).Seconds()

	if this.accLog != nil {
		ainfo := make(map[string]interface{})
		ainfo["url"] = qurl
		ainfo["status"] = hresp.StatusCode
		if len(content) < 100 {
			ainfo["content"] = content
		} else {
			ainfo["content"] = string(content[:100])
		}
		ar := acclog.NewCommonLog(ainfo, te.Sub(ts).Seconds())
		this.accLog.Write(this.accName, ar)
	}
	return m, nil
}

// httpclient.exec(req:table) resp
type GOF_httpclient_exec struct {
	m *httpClientModule
}

func (this *GOF_httpclient_exec) Exec(vm *golua.VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	reqd, err1 := vm.API_pop1X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vreqd := vm.API_table(reqd)
	if vreqd == nil {
		return 0, fmt.Errorf("req param invalid(%T)", reqd)
	}
	m := golua.GoData(vreqd).(map[string]interface{})
	req := new(httpclientRequest)
	if !valutil.ToBean(m, req) {
		return 0, fmt.Errorf("httpclientRequest invalid(%v)", m)
	}
	rm, err2 := this.m.doExecute(vm, req)
	if err2 != nil {
		return 0, err2
	}
	vm.API_push(vm.API_table(rm))
	return 1, nil
}

func (this *GOF_httpclient_exec) IsNative() bool {
	return true
}

func (this *GOF_httpclient_exec) String() string {
	return "GoFunc<httpclient.exec>"
}

// httpclient.getContent(url:string[, tmMS:int]) string
type GOF_httpclient_getContent struct {
	m *httpClientModule
}

func (this *GOF_httpclient_getContent) Exec(vm *golua.VM, self interface{}) (int, error) {
	err0 := vm.API_checkStack(1)
	if err0 != nil {
		return 0, err0
	}
	url, tm, err1 := vm.API_pop2X(-1, true)
	if err1 != nil {
		return 0, err1
	}
	vurl := valutil.ToString(url, "")
	if vurl == "" {
		return 0, fmt.Errorf("url invalid(%v)", url)
	}
	req := new(httpclientRequest)
	req.URL = vurl
	req.TimeoutMS = valutil.ToInt(tm, 0)
	rm, err2 := this.m.doExecute(vm, req)
	if err2 != nil {
		return 0, err2
	}
	var content interface{}
	st := valutil.ToInt(rm["Status"], 0)
	if st == http.StatusOK {
		content = rm["Content"]
	} else {
		content = ""
	}
	vm.API_push(content)
	return 1, nil
}

func (this *GOF_httpclient_getContent) IsNative() bool {
	return true
}

func (this *GOF_httpclient_getContent) String() string {
	return "GoFunc<httpclient.getContent>"
}
