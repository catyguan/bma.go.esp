package esnp

import (
	"bmautil/valutil"
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"time"
)

type URL struct {
	Uri   []string
	Data  *url.URL
	Query url.Values
}

func ParseURL(s string) (*URL, error) {
	v, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	r := new(URL)
	r.Data = v
	uri := v.RequestURI()
	slist := strings.SplitN(uri, "/", 3)
	if len(slist) > 1 {
		r.Uri = slist[1:]
	} else {
		r.Uri = make([]string, 0)
	}
	r.Query = r.Data.Query()
	return r, nil
}

func (this *URL) GetHost() string {
	if strings.ToLower(this.Data.Host) != "unknow" {
		return this.Data.Host
	}
	return ""
}

func (this *URL) GetService() string {
	if len(this.Uri) > 0 {
		if strings.ToLower(this.Uri[0]) != "unknow" {
			return this.Uri[0]
		}
	}
	return ""
}

func (this *URL) GetOp() string {
	if len(this.Uri) > 1 {
		if strings.ToLower(this.Uri[1]) != "unknow" {
			return this.Uri[1]
		}
	}
	return ""
}

func (this *URL) GetGroup() string {
	qs := this.Query
	v := qs.Get("g")
	if v != "" {
		return v
	}
	v = qs.Get("a_50")
	if v != "" {
		return v
	}
	return ""
}

func (this *URL) GetObject() string {
	qs := this.Query
	v := qs.Get("o")
	if v != "" {
		return v
	}
	v = qs.Get("a_10")
	if v != "" {
		return v
	}
	return ""
}

func (this *URL) GetTimeout(d time.Duration) time.Duration {
	qs := this.Query
	v := valutil.ToInt(qs.Get("to"), 0)
	if v <= 0 {
		return d
	}
	return time.Duration(v) * time.Millisecond
}

func (this *URL) BindAddress(r *Address, bhost bool) error {
	var v string
	v = this.GetHost()
	if bhost && v != "" {
		r.SetHost(v)
	}
	v = this.GetService()
	if v != "" {
		r.SetService(v)
	}
	v = this.GetOp()
	if v != "" {
		r.SetOp(v)
	}
	qs := this.Query
	for k, _ := range qs {
		qv := qs.Get(k)
		if k == "o" && len(qv) > 0 {
			r.SetObject(qv)
			continue
		}
		if k == "g" {
			r.SetGroup(qv)
		}
		if strings.HasPrefix(k, "a_") {
			nk := valutil.ToInt(k[2:], 0)
			if nk > 0 {
				r.Set(nk, qv)
			}
		}
	}
	return nil
}

func (this *URL) BindMessage(msg *Message, bhost bool) error {
	addr := msg.GetAddress()
	return this.BindAddress(addr, bhost)
}

func (this *Address) ToURL() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("esnp://")
	if host := this.GetHost(); host != "" {
		buf.WriteString(host)
	} else {
		buf.WriteString("unknow")
	}
	buf.WriteString("/")
	if service := this.GetService(); service != "" {
		buf.WriteString(service)
	} else {
		buf.WriteString("unknow")
	}
	buf.WriteString("/")
	if op := this.GetOp(); op != "" {
		buf.WriteString(op)
	} else {
		buf.WriteString("unknow")
	}

	b := false
	anns := this.Annotations()
	for _, ann := range anns {
		switch ann {
		case ADDRESS_HOST, ADDRESS_OP, ADDRESS_SERVICE:
			continue
		}
		n := ""
		v := this.Get(ann)
		switch ann {
		case ADDRESS_OBJECT:
			n = "o"
		case ADDRESS_GROUP:
			n = "g"
		default:
			n = fmt.Sprintf("a%d", ann)
		}
		if !b {
			b = true
			buf.WriteString("?")
		} else {
			buf.WriteString("&")
		}
		buf.WriteString(n)
		buf.WriteString("=")
		buf.WriteString(url.QueryEscape(v))
	}
	return buf.String()
}
