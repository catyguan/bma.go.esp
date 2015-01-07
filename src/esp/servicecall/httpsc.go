package servicecall

import (
	"bmautil/httputil"
	"bmautil/valutil"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"logger"
	"net/http"
	"net/url"
	"objfac"
	"strings"
	"time"
)

type HttpServiceCaller struct {
	name string
	httpConfig
}

func (this *HttpServiceCaller) SetName(n string) {
	this.name = n
}

func (this *HttpServiceCaller) Start() error {
	return nil
}

func (this *HttpServiceCaller) Stop() {
}

func (this *HttpServiceCaller) Call(serviceName, method string, params map[string]interface{}, deadline time.Time) (interface{}, error) {
	bs, err0 := json.Marshal(params)
	if err0 != nil {
		return nil, err0
	}
	qurl := this.URL
	data := make(url.Values)
	data.Add("m", method)
	data.Add("p", string(bs))
	body := strings.NewReader(data.Encode())

	hreq, err2 := http.NewRequest("POST", qurl, body)
	if err2 != nil {
		return nil, err2
	}
	hreq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if this.Host != "" {
		hreq.Header.Set("Host", this.Host)
	}

	tm := this.TimeoutMS
	if tm <= 0 {
		tm = 5000
	}
	tmd := time.Duration(tm) * time.Millisecond
	if !deadline.IsZero() {
		timeout := deadline.Sub(time.Now())
		if timeout < tmd {
			tmd = timeout
		}
	}
	client := httputil.NewHttpClient(tmd)

	ts := time.Now()
	hresp, err3 := client.Do(hreq)
	te := time.Now()
	if err3 != nil {
		logger.Debug(tag, "[%s:%s] http '%s'(%f) fail '%s'", serviceName, method, qurl, te.Sub(ts).Seconds(), err3)
		return nil, err3
	}
	logger.Debug(tag, "[%s：%s] http '%s'(%f) end '%d'", serviceName, method, qurl, te.Sub(ts).Seconds(), hresp.StatusCode)
	defer hresp.Body.Close()
	if hresp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("'%s' invalid http status(%d)", qurl, hresp.StatusCode)
	}

	respBody, err4 := ioutil.ReadAll(hresp.Body)
	if err4 != nil {
		return nil, err4
	}
	var r interface{}
	err5 := json.Unmarshal(respBody, &r)
	if err5 != nil {
		return nil, err5
	}
	return r, nil
}

type httpConfig struct {
	URL       string
	Host      string
	TimeoutMS int
}

type HttpServiceCallerFactory int

func (o HttpServiceCallerFactory) Valid(cfg map[string]interface{}, ofp objfac.ObjectFactoryProvider) error {
	var co httpConfig
	if valutil.ToBean(cfg, &co) {
		if co.URL == "" {
			return fmt.Errorf("URL empty")
		}
		return nil
	}
	return fmt.Errorf("invalid HttpServiceCaller config")
}

func (o HttpServiceCallerFactory) Compare(cfg map[string]interface{}, old map[string]interface{}, ofp objfac.ObjectFactoryProvider) (same bool) {
	var co, oo httpConfig
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

func (o HttpServiceCallerFactory) Create(cfg map[string]interface{}, ofp objfac.ObjectFactoryProvider) (interface{}, error) {
	err := o.Valid(cfg, ofp)
	if err != nil {
		return nil, err
	}
	var co httpConfig
	valutil.ToBean(cfg, &co)
	r := new(HttpServiceCaller)
	r.httpConfig = co
	return ServiceCaller(r), nil
}
