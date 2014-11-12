package n2n

import (
	"esp/cluster/nodebase"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestService(t *testing.T) {
	time.AfterFunc(5*time.Second, func() {
		fmt.Println("os exit!!!")
		os.Exit(-1)
	})

	nodebase.Id = nodebase.NodeId(100)
	nodebase.Name = "testcase"

	cfg := new(ConfigInfo)
	cfg.Host = "_:9999"
	cfg.Remotes = make(MapOfRemoteConfigInfo)
	if true {
		rc := new(RemoteConfigInfo)
		rc.Host = "127.0.0.1:1090"
		rc.Code = "123"
		cfg.Remotes["coo"] = rc
	}
	cfg.Valid()

	s := NewService(8)
	s.InitConfig(cfg)
	s.Start()
	s.Run()
	time.Sleep(1 * time.Second)
	s.Stop()
	time.Sleep(100 * time.Millisecond)
}
