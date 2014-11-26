package espnetss

import (
	"bytes"
	"esp/espnet/espsocket"
	"fmt"
	"strings"
)

const (
	tag = "espnetsdk"
)

var (
	gHandlers map[string]LoginHandler
)

func init() {
	gHandlers = make(map[string]LoginHandler)
	RegisterLoginHandler("none", noneLogin)
}

func RegisterLoginHandler(typ string, lg LoginHandler) {
	gHandlers[typ] = lg
}

func GetLoginHandler(typ string) LoginHandler {
	if lh, ok := gHandlers[typ]; ok {
		return lh
	}
	return nil
}

type LoginHandler func(sock *espsocket.Socket, user string, cert string) (bool, error)

type Config struct {
	Host        string
	User        string
	LoginType   string
	Certificate string
	PoolSize    int
	PreConns    int
}

func (this *Config) Parse(s string) {
	this.Host, this.User, this.LoginType, this.Certificate = Split(s)
}

func (this *Config) Valid() error {
	if this.Host == "" {
		return fmt.Errorf("Host empty")
	}
	if this.LoginType != "" {
		li := GetLoginHandler(this.LoginType)
		if li == nil {
			return fmt.Errorf("LoginType(%s) miss", this.LoginType)
		}
	}
	return nil
}

func (this *Config) Compare(o *Config) bool {
	if this.Host != o.Host {
		return false
	}
	if this.User != o.User {
		return false
	}
	if this.LoginType != o.LoginType {
		return false
	}
	if this.Certificate != o.Certificate {
		return false
	}
	if this.PoolSize != o.PoolSize {
		return false
	}
	if this.PreConns != o.PreConns {
		return false
	}
	return true
}

func (this *Config) Key() string {
	return Make(this.Host, this.User, this.LoginType, this.Certificate)
}

func Make(host, user, lt, cert string) string {
	buf := bytes.NewBuffer([]byte{})
	if user != "" {
		buf.WriteString(user)
	}
	if lt != "" {
		if buf.Len() > 0 {
			buf.WriteByte(':')
		}
		buf.WriteString(lt)
		if cert != "" {
			buf.WriteByte(',')
			buf.WriteString(cert)
		}
	}
	if host != "" {
		if buf.Len() > 0 {
			buf.WriteByte('@')
		}
		buf.WriteString(host)
	}
	return buf.String()
}

func Split(netsource string) (host, user, lt, cert string) {
	str := netsource
	i1 := strings.LastIndex(str, "@")
	if i1 != -1 {
		host = str[i1+1:]
		str = str[:i1]
	} else {
		host = str
		return
	}
	i2 := strings.Index(str, ":")
	if i2 != -1 {
		user = str[:i2]
		str = str[i2+1:]
	} else {
		user = str
		return
	}
	i3 := strings.Index(str, ",")
	if i3 != -1 {
		lt = str[:i3]
		cert = str[i3+1:]
	} else {
		lt = str
	}
	return host, user, lt, cert
}
