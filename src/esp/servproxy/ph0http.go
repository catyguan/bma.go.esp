package servproxy

import (
	"fmt"
	"net"
)

type HttpProxyReq struct {
	s    *Service
	port *PortObj
	conn net.Conn
}

type HttpProxyHandler int

func init() {
	AddProxyHandler("http", HttpProxyHandler(0))
}

func (this HttpProxyHandler) Handle(s *Service, port *PortObj, conn net.Conn) {
	req := new(HttpProxyReq)
	req.s = s
	conn.Close()
}

func (this HttpProxyHandler) AnswerError(port *PortObj, req interface{}, err error) error {
	return err
}

func (this HttpProxyHandler) Valid(cfg *RemoteConfigInfo) error {
	if cfg.Host == "" {
		return fmt.Errorf("Host invalid")
	}
	return nil
}

func (this HttpProxyHandler) Compare(cfg *RemoteConfigInfo, old *RemoteConfigInfo) bool {
	return true
}

func (this HttpProxyHandler) Start(o *RemoteObj) error {
	o.Data = nil
	return nil
}

func (this HttpProxyHandler) Stop(o *RemoteObj) error {
	// http.Serve(l, handler)
	return nil
}

func (this HttpProxyHandler) Forward(port *PortObj, req interface{}, remote *RemoteObj) error {
	return nil
}
