package simplelogin

import (
	"acl"
	"bmautil/lru"
	"boot"
	"fmt"
	"io/ioutil"
	"logger"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"uuid"
)

const (
	tag = "simplelogin"
)

type SimpleLoginMux struct {
	file      string
	sessionId string
	h         http.Handler
	lock      sync.RWMutex
	users     *lru.Cache
}

func NewSimpleLoginMux(pwfile string, sessionId string, cacheSize int32, h http.Handler) *SimpleLoginMux {
	r := new(SimpleLoginMux)
	r.file = pwfile
	if sessionId == "" {
		sessionId = "SLM_SESSION_ID"
	}
	r.sessionId = sessionId
	r.h = h
	r.users = lru.NewCache(cacheSize)
	return r
}

func (this *SimpleLoginMux) BindHandler(h http.Handler) {
	this.h = h
}

func (this *SimpleLoginMux) UserProvider(r *http.Request) (*acl.User, error) {
	ck, err := r.Cookie(this.sessionId)
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, nil
		}
		return nil, err
	}
	sid := ck.Value
	// logger.Debug(tag, "session id - %s", sid)
	this.lock.RLock()
	defer this.lock.RUnlock()
	v, _ := this.users.Get(sid)
	if v == nil {
		return nil, nil
	}
	return v.(*acl.User), nil
}

func (this *SimpleLoginMux) createSessionId() string {
	var s string
	u, err := uuid.NewV4()
	if err != nil {
		s = fmt.Sprintf("%d", time.Now().Nanosecond())
	} else {
		s = u.String()
	}
	return s
}

func (this *SimpleLoginMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, err0 := this.UserProvider(r)
	if err0 != nil {
		logger.Error(tag, "find user fail - %s", err0)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if user == nil {
		logger.Debug(tag, "anonymous access %s", r.URL.Path)
		if strings.HasSuffix(r.URL.Path, "favicon.ico") {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		msg := ""
		if true {
			r.ParseForm()
			un := r.PostFormValue("user")
			pw := r.PostFormValue("pw")
			if un != "" {
				ip, _, _ := net.SplitHostPort(r.RemoteAddr)
				user, err0 = this.doLogin(un, pw, ip)
				if err0 != nil {
					logger.Error(tag, "doLogin fail - %s", err0)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				if user == nil {
					msg = "<font color=red>login fail</font><br>"
				} else {
					sid := this.createSessionId()

					this.lock.Lock()
					defer this.lock.Unlock()
					this.users.Put(sid, user)

					ck := new(http.Cookie)
					ck.Name = this.sessionId
					ck.Value = sid
					ck.Path = "/"
					http.SetCookie(w, ck)
					ref := r.Header.Get("Referer")
					w.Header().Set("Location", ref)
					w.WriteHeader(http.StatusMovedPermanently)
					return
				}
			}
		}
		if user == nil {
			logger.Debug(tag, "send login form")
			this.showLoginForm(w, r, msg)
			return
		}
	}
	this.h.ServeHTTP(w, r)
}

var loginForm = `<html>
<head>
<title>login</title>
</head>
<body>
%s
<form method="POST">
User:<input type="text" name="user"><br>
Pass:<input type="password" name="pw"><br>
<input type="submit" value="Login">
</form>
</body>
</html>
`

func (this *SimpleLoginMux) showLoginForm(w http.ResponseWriter, r *http.Request, msg string) {
	fmt.Fprintf(w, loginForm, msg)
}

func (this *SimpleLoginMux) doLogin(user, pass string, ip string) (*acl.User, error) {
	fn := this.file
	if fn == "" {
		fn = filepath.Join(boot.WorkDir, "slusers.dat")
	}
	bs, err0 := ioutil.ReadFile(fn)
	if err0 != nil {
		return nil, err0
	}
	list := strings.Split(string(bs), "\n")
	for _, s := range list {
		plist := strings.Split(strings.TrimSpace(s), ":")
		if plist[0] == user {
			if len(plist) > 1 && plist[1] == pass {
				var g []string
				if len(plist) > 2 {
					sg := plist[2]
					g = strings.Split(sg, ",")
				}
				r := acl.NewUser(user, ip, g)
				logger.Debug(tag, "'%s' login -> %s, %v", user, r, g)
				return r, nil
			} else {
				logger.Debug(tag, "'%s' pass invalid", user)
				return nil, nil
			}
		}
	}
	logger.Debug(tag, "'%s' user invalid", user)
	return nil, nil
}
