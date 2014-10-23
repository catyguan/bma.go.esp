package main

import (
	"boot"
	"bytes"
	"crypto/md5"
	"fileloader"
	"fmt"
	"logger"
	"mime"
	"net"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

type Service struct {
	name   string
	config *configInfo
	fl     fileloader.FileLoader
}

func (this *Service) createCode(ip, module, file, version string) string {
	key := this.config.Key
	buf := bytes.NewBuffer(make([]byte, 0, 64))
	buf.WriteString(ip)
	buf.WriteString("`")
	buf.WriteString(module)
	buf.WriteString("`")
	// buf.WriteString(file)
	// buf.WriteString("`")
	// buf.WriteString(version)
	// buf.WriteString("`")
	buf.WriteString(key)
	h := md5.New()
	h.Write(buf.Bytes())
	r := fmt.Sprintf("%x", h.Sum(nil))
	// if len(r) > 32 {
	// 	r = r[len(r)-32:]
	// }
	return r
}

func (this *Service) InvokeCreate(w http.ResponseWriter, req *http.Request) {
	i := req.FormValue("i")
	module := req.FormValue("m")
	file := req.FormValue("f")
	version := req.FormValue("v")

	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	if !boot.DevMode && ip != this.config.AdminIp {
		http.Error(w, fmt.Sprintf("invalid admin ip(%s)", ip), http.StatusForbidden)
		return
	}
	if i == "" {
		i = ip
	}
	code := this.createCode(i, module, file, version)

	str := req.RequestURI
	idx := strings.LastIndex(str, "/")
	if idx == -1 {
		str = "query"
	} else {
		str = str[:idx+1] + "query"
	}

	buf := bytes.NewBuffer(make([]byte, 0, 256))
	buf.WriteString("http://")
	buf.WriteString(req.Host)
	buf.WriteString(str)
	buf.WriteString("?")
	buf.WriteString("m=")
	buf.WriteString(module)
	buf.WriteString("&f=")
	if file == "" {
		file = fileloader.VAR_F
	}
	buf.WriteString(file)
	buf.WriteString("&v=")
	buf.WriteString(version)
	buf.WriteString("&c=")
	buf.WriteString(code)

	w.Write(buf.Bytes())
}

func (this *Service) InvokeFL(w http.ResponseWriter, req *http.Request) {
	ip, _, _ := net.SplitHostPort(req.RemoteAddr)
	module := req.FormValue("m")
	file := req.FormValue("f")
	version := req.FormValue("v")
	code := req.FormValue("c")
	mycode := this.createCode(ip, module, file, version)

	if code != mycode {
		logger.Warn(tag, "invalid code(i=%s,m=%s,f=%s,v=%s) in:%s my:%s", ip, module, file, version, code, mycode)
		http.Error(w, fmt.Sprintf("invalid code for '%s'", ip), http.StatusForbidden)
		return
	}

	if version == "" {
		version = "default"
	}
	if strings.HasPrefix(file, "/") {
		file = file[1:]
	}
	fn := fmt.Sprintf("%s:%s/%s", module, version, file)
	logger.Debug(tag, "query(%s, %s, %s) -> %s", module, version, file, fn)

	bs, err := this.fl.Load(fn)
	if err != nil {
		logger.Debug(tag, "query(%s) fail - %s", fn, err)
		http.Error(w, fmt.Sprintf("access %s:%s#%s fail", module, file, version), http.StatusBadRequest)
		return
	}
	if bs == nil {
		logger.Debug(tag, "query(%s) miss", fn)
		http.NotFound(w, req)
		return
	}
	l := int64(len(bs))
	logger.Debug(tag, "load %s -> ok", fn)
	ctype := mime.TypeByExtension(filepath.Ext(file))
	if ctype == "" {
		ctype = http.DetectContentType(bs)
	}
	w.Header().Set("Content-Type", ctype)
	w.Header().Set("Content-Length", strconv.FormatInt(l, 10))
	w.WriteHeader(http.StatusOK)

	w.Write(bs)
}
