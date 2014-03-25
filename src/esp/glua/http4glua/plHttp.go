package http4glua

import (
	"bmautil/valutil"
	"esp/glua"
	"fmt"
	"io"
	"io/ioutil"
	"logger"
	"lua51"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	tag = "http4glua"
)

type Request struct {
	URL       string
	Headers   map[string]string
	Post      bool
	Data      map[string]string
	ResultKey string
}

func (this *Request) Valid() error {
	if this.URL == "" {
		return fmt.Errorf("url empty")
	}
	return nil
}

type PluginHttp struct {
}

func (tshi *PluginHttp) Name() string {
	return "http"
}

func (this *PluginHttp) OnInitLua(l *lua51.State) error {
	return nil
}

func (this *PluginHttp) OnCloseLua(l *lua51.State) {
}

func (this *PluginHttp) Execute(task *glua.PluginTask) error {
	req := new(Request)
	if task.Request != nil {
		valutil.ToBean(task.Request, req)
	}
	err := req.Valid()
	if err != nil {
		return err
	}
	go func() {
		err := this.doExecute(task, req)
		if err != nil {
			task.Callback(this.Name(), nil, err)
		}
	}()
	return nil
}

func (this *PluginHttp) doExecute(task *glua.PluginTask, req *Request) error {
	var body io.Reader
	method := "GET"
	if req.Post {
		method = "POST"
		data := make(url.Values)
		for k, v := range req.Data {
			data.Add(k, v)
		}
		body = strings.NewReader(data.Encode())
	}
	hreq, err2 := http.NewRequest(method, req.URL, body)
	if err2 != nil {
		return err2
	}
	if req.Post {
		hreq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	for k, v := range req.Headers {
		hreq.Header.Set(k, v)
	}
	client := http.DefaultClient
	logger.Debug(tag, "[%s] http '%s' start", task.Context, req.URL)
	ts := time.Now()
	hresp, err3 := client.Do(hreq)
	te := time.Now()
	if err3 != nil {
		logger.Debug(tag, "[%s] http '%s' fail '%s'", task.Context, req.URL, err3)
		return err3
	}
	logger.Debug(tag, "[%s] http '%s' end '%d'", task.Context, req.URL, hresp.StatusCode)
	defer hresp.Body.Close()
	respBody, err4 := ioutil.ReadAll(hresp.Body)
	if err4 != nil {
		return err4
	}
	m := make(map[string]interface{})
	m["StatusCode"] = hresp.StatusCode
	hs := make(map[string]string)
	for k, _ := range hresp.Header {
		v := hresp.Header.Get(k)
		hs[k] = v
	}
	m["Header"] = hs
	m["Content"] = string(respBody)
	m["Time"] = te.Sub(ts).Seconds()

	task.Callback(this.Name(), func(ctx *glua.Context) {
		rk := req.ResultKey
		if rk == "" {
			rk = "http"
		}
		ctx.Result[rk] = m
	}, nil)
	return nil
}
