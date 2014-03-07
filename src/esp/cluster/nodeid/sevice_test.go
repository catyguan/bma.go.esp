package nodeid

import (
	"boot"
	"testing"
)

func TestNodeId(t *testing.T) {
	cfile := "../../../../bin/config/xmem-config.json"

	nodeId := NewService("espnode")
	boot.QuickDefine(nodeId, "", true)

	f1 := func() {

	}
	if f1 != nil {
	}

	funl := []func(){
		f1,
	}

	boot.TestGo(cfile, 1, funl)
}
