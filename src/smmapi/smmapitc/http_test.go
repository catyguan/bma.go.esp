package smmapitc

import "testing"

const (
	NODE_HTTP_URL  = "http://127.0.0.1:1081/smm.api/invoke"
	NODE_HTTP_CODE = "1"
)

func TestReload(t *testing.T) {
	httpInvoke("go.server", "boot.reload", nil)
}

func TestShutdown(t *testing.T) {
	httpInvoke("go.server", "boot.shutdown", nil)
}

func TestProfHeap(t *testing.T) {
	httpInvoke("go.pprof", "pprof.heap", nil)
}
