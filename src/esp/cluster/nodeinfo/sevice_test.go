package nodeinfo

import (
	"boot"
	"testing"
)

func TestNodeId(t *testing.T) {
	cfile := "../../../../bin/config/xmem-config.json"

	nodeInfo := NewService("espnode")
	boot.Add(nodeInfo, "", true)

	f1 := func() {

	}
	if f1 != nil {
	}

	funl := []func(){
		f1,
	}

	boot.TestGo(cfile, 1, funl)
}
