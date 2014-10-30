package aclmux

import (
	"acl"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type AclServerMux struct {
	name string
	h    http.Handler
}

func NewAclServerMux(n string, h http.Handler) *AclServerMux {
	r := new(AclServerMux)
	r.name = n
	r.h = h
	return r
}

func (this *AclServerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	user := acl.NewUser("anonymous", ip, nil)

	path := r.URL.Path
	if strings.HasPrefix(path, "/") {
		path = this.name + path
	} else {
		path = this.name + "/" + path
	}
	ps := strings.Split(path, "/")

	err := acl.Assert(user, ps, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("IP(%s) FORBIDDEN", ip), http.StatusForbidden)
		return
	}
	this.h.ServeHTTP(w, r)
}
