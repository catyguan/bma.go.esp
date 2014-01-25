package nodeid

import (
	"boot"
	"esp/sqlite"
	"fmt"
	"testing"
)

func TestNodeId(t *testing.T) {
	cfile := "../../../bin/config/xmem-config.json"

	sqliteServer := sqlite.NewSqliteServer("sqliteServer")
	sqliteServer.DefaultBoot()

	nodeId := NewService("espnode", sqliteServer)
	boot.QuickDefine(nodeId, "", true)

	f1 := func() {
		fmt.Println("nodeId", nodeId.GetAndListen("test", func(nid uint64) {
			fmt.Println("newNodeId", nid)
		}))
		nodeId.SetId(123)
	}
	if f1 != nil {
	}

	funl := []func(){
		f1,
	}

	boot.TestGo(cfile, 1, funl)
}
