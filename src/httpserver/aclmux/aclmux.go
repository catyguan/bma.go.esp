package aclmux

import (
	"acl"
	"fmt"
	"logger"
	"net"
	"net/http"
	"strings"
)

type AclServerMux struct {
	name string
	h    http.Handler
	up   UserProvider
}

func NewAclServerMux(n string, h http.Handler, up UserProvider) *AclServerMux {
	r := new(AclServerMux)
	r.name = n
	r.h = h
	r.up = up
	return r
}

func (this *AclServerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var user *acl.User
	if this.up != nil {
		var err error
		user, err = this.up(r)
		if err != nil {
			logger.Error("aclmux", "find user fail - %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	if user == nil {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		user = acl.NewUser("anonymous", ip, nil)
	}

	path := r.URL.Path
	if strings.HasPrefix(path, "/") {
		path = this.name + path
	} else {
		path = this.name + "/" + path
	}
	ps := strings.Split(path, "/")

	err := acl.Assert(user, ps, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("USER(%s) FORBIDDEN", user), http.StatusForbidden)
		return
	}
	this.h.ServeHTTP(w, r)
}
