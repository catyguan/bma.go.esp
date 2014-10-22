package http4fileloader

import (
	"bmautil/httputil"
	"bmautil/valutil"
	"fileloader"
	"fmt"
	"io/ioutil"
	"logger"
	"net/http"
	"strings"
	"time"
)

const (
	tag = "httpfileloader"
)

func init() {
	fileloader.AddFileLoaderFactory("http", HttpFileLoaderFactory)
}

type HttpFileLoader struct {
	config
}

func (this *HttpFileLoader) Load(script string) ([]byte, error) {
	module, n := fileloader.SplitModuleScript(script)

	method := "GET"
	qurl := this.URL
	if strings.Contains(qurl, fileloader.VAR_M) {
		strings.Replace(qurl, fileloader.VAR_M, module, -1)
	} else {
		qurl += module + "/"
	}
	if strings.Contains(qurl, fileloader.VAR_F) {
		strings.Replace(qurl, fileloader.VAR_F, n, -1)
	} else {
		qurl += n
	}
	hreq, err2 := http.NewRequest(method, qurl, nil)
	if err2 != nil {
		return nil, err2
	}
	if this.Host != "" {
		hreq.Header.Set("Host", this.Host)
	}

	tm := this.TimeoutMS
	if tm <= 0 {
		tm = 5000
	}
	client := httputil.NewHttpClient(time.Millisecond * time.Duration(tm))

	ts := time.Now()
	hresp, err3 := client.Do(hreq)
	te := time.Now()
	if err3 != nil {
		logger.Debug(tag, "[%s] http '%s'(%f) fail '%s'", script, qurl, te.Sub(ts).Seconds(), err3)
		return nil, err3
	}
	logger.Debug(tag, "[%s] http '%s'(%f) end '%d'", script, qurl, te.Sub(ts).Seconds(), hresp.StatusCode)
	defer hresp.Body.Close()
	if hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid http status(%d)", hresp.StatusCode)
	}

	respBody, err4 := ioutil.ReadAll(hresp.Body)
	if err4 != nil {
		return nil, err4
	}
	return respBody, nil
}

type config struct {
	URL       string
	Host      string
	TimeoutMS int
}

type httpFileLoaderFactory int

const (
	HttpFileLoaderFactory = httpFileLoaderFactory(0)
)

func (this httpFileLoaderFactory) Valid(cfg map[string]interface{}) error {
	var co config
	if valutil.ToBean(cfg, &co) {
		if co.URL == "" {
			return fmt.Errorf("URL empty")
		}
		return nil
	}
	return fmt.Errorf("invalid HttpFileLoader config")
}

func (this httpFileLoaderFactory) Compare(cfg map[string]interface{}, old map[string]interface{}) bool {
	var co, oo config
	if !valutil.ToBean(cfg, &co) {
		return false
	}
	if !valutil.ToBean(old, &oo) {
		return false
	}
	if co.URL != oo.URL {
		return false
	}
	if co.Host != oo.Host {
		return false
	}
	if co.TimeoutMS != oo.TimeoutMS {
		return false
	}
	return true
}

func (this httpFileLoaderFactory) Create(cfg map[string]interface{}) (fileloader.FileLoader, error) {
	err := this.Valid(cfg)
	if err != nil {
		return nil, err
	}
	var co config
	valutil.ToBean(cfg, &co)
	r := new(HttpFileLoader)
	r.config = co
	return r, nil
}
