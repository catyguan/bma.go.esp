package n2n

import (
	"esp/cluster/nodeinfo"
	"esp/espnet/esnp"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestConnector(t *testing.T) {
	time.AfterFunc(5*time.Second, func() {
		fmt.Println("os exit!!!")
		os.Exit(-1)
	})

	ninfo := nodeinfo.NewService("nodeInfo")

	s := NewService("test", ninfo)
	ctor := new(connector)
	if true {
		url, err := esnp.ParseURL("esnp://127.0.0.1:1080/")
		if err != nil {
			t.Error(err)
			return
		}
		ctor.InitConnector(s, "tc", url)
	}
	time.Sleep(1 * time.Second)
	ctor.Close()
}
