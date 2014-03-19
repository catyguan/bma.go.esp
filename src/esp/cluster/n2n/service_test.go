package n2n

import (
	"bmautil/socket"
	"boot"
	"esp/cluster/nodeinfo"
	"esp/espnet/espservice"
	"testing"
)

func TestService1(t *testing.T) {
	cfgFile := "service_test1.json"
	doServiceTest(cfgFile)
}

func TestService2(t *testing.T) {
	cfgFile := "service_test2.json"
	doServiceTest(cfgFile)
}

func doServiceTest(cfgFile string) {
	ninfo := nodeinfo.NewService("nodeInfo")
	boot.Add(ninfo, "", false)

	s := NewService("n2nService", ninfo)
	boot.Add(s, "", false)

	mux := espservice.NewServiceMux(nil, nil)
	mux.AddServiceHandler("n2n", s.Serve)
	goService := espservice.NewGoService("goService", mux.Serve)

	lisPoint := socket.NewListenPoint("servicePoint", nil, goService.AcceptESP)
	boot.Add(lisPoint, "", false)

	boot.TestGo(cfgFile, 15, nil)
}
