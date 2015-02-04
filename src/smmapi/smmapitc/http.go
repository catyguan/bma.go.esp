package smmapitc

import (
	"bmautil/httputil"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"logger"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	tag = "smmapitc"
)

func httpInvoke(id, aid string, param map[string]interface{}) {
	resp, err := func() (map[string]interface{}, error) {
		var body io.Reader
		post := false
		method := "GET"
		qurl := NODE_HTTP_URL
		data := make(url.Values)
		data.Add("id", id)
		data.Add("aid", aid)
		strparam := ""
		if len(param) > 0 {
			buf, err2 := json.Marshal(data)
			if err2 != nil {
				return nil, err2
			}
			strparam = string(buf)
			data.Add("param", strparam)
			post = true
		}
		if NODE_HTTP_CODE != "" {
			tmp := fmt.Sprintf("%s/%s/%s/%s", id, aid, strparam, NODE_HTTP_CODE)
			h := md5.New()
			h.Write([]byte(tmp))
			mycode := fmt.Sprintf("%x", h.Sum(nil))
			data.Add("code", mycode)
		}

		if post {
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
		if post {
			hreq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		// for k, v := range req.Headers {
		// 	hreq.Header.Set(k, v)
		// }
		client := httputil.NewHttpClient(time.Millisecond * time.Duration(5000))
		logger.Debug(tag, "http '%s' start", qurl)

		ts := time.Now()
		hresp, err3 := client.Do(hreq)
		te := time.Now()
		if err3 != nil {
			logger.Debug(tag, "http '%s' fail '%s'", qurl, err3)
			return nil, err3
		}
		logger.Debug(tag, "http '%s' end '%d'", qurl, hresp.StatusCode)
		defer hresp.Body.Close()
		respBody, err4 := ioutil.ReadAll(hresp.Body)
		if err4 != nil {
			return nil, err4
		}
		m := make(map[string]interface{})
		m["Status"] = hresp.StatusCode
		content := string(respBody)
		m["Content"] = content
		m["Time"] = te.Sub(ts).Seconds()
		return m, nil
	}()

	if err != nil {
		fmt.Println("error", err)
	} else {
		fmt.Println(resp)
	}
}
