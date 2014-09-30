package http4glua

import (
	"bmautil/valutil"
	"context"
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
	Data      map[string]interface{}
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
	ctx := task.Context

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
		return err2
	}
	if req.Post {
		hreq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	for k, v := range req.Headers {
		hreq.Header.Set(k, v)
	}
	logstr := ""
	if logger.EnableDebug(tag) {
		logstr = glua.GLuaContext.String(task.Context)
	}
	client := http.DefaultClient
	logger.Debug(tag, "[%s] http '%s' start", logstr, qurl)
	if glua.GLuaContext.HasAccessLog(ctx) {
		ainfo := make(map[string]interface{})
		ainfo["url"] = qurl
		glua.GLuaContext.DoAccessLog(ctx, "http:start", nil)
	}

	ts := time.Now()
	hresp, err3 := client.Do(hreq)
	te := time.Now()
	if err3 != nil {
		logger.Debug(tag, "[%s] http '%s' fail '%s'", logstr, qurl, err3)
		if glua.GLuaContext.HasAccessLog(ctx) {
			ainfo := make(map[string]interface{})
			ainfo["url"] = qurl
			ainfo["error"] = err3.Error()
			glua.GLuaContext.DoAccessLog(ctx, "http:end", ainfo)
		}
		return err3
	}
	logger.Debug(tag, "[%s] http '%s' end '%d'", logstr, qurl, hresp.StatusCode)
	defer hresp.Body.Close()
	respBody, err4 := ioutil.ReadAll(hresp.Body)
	if err4 != nil {
		if glua.GLuaContext.HasAccessLog(ctx) {
			ainfo := make(map[string]interface{})
			ainfo["url"] = qurl
			ainfo["error"] = err4.Error()
			glua.GLuaContext.DoAccessLog(ctx, "http:end", ainfo)
		}
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

	glua.GLuaContext.DoAccessLog(ctx, "http:end", nil)

	task.Callback(this.Name(), func(ctx context.Context) {
		rk := req.ResultKey
		if rk == "" {
			rk = "http"
		}
		rs := glua.GLuaContext.GetResult(ctx)
		rs[rk] = m
	}, nil)
	return nil
}
