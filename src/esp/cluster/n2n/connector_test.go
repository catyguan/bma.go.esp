package n2n

import (
	"esp/cluster/nodebase"
	"fmt"
	"os"
	"testing"
	"time"
)

func T2estConnector(t *testing.T) {
	time.AfterFunc(5*time.Second, func() {
		fmt.Println("os exit!!!")
		os.Exit(-1)
	})

	nodebase.Id = nodebase.NodeId(100)
	nodebase.Name = "testcase"

	s := NewService(8)
	ctor := new(connector)
	if true {
		ctor.InitConnector(s, "tc", "127.0.0.1:1090", "")
	}
	time.Sleep(1 * time.Second)
	ctor.Close()
}
