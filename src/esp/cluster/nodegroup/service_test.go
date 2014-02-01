package nodegroup

import (
	"boot"
	"esp/cluster/election"
	"esp/cluster/nodeid"
	"esp/espnet"
	"esp/sqlite"
	"fmt"
	"logger"
	"testing"
)

func func4tester(s *Service) func() {
	return func() {
		cfg := new(NodeGroupConfig)
		cfg.IdleCheckMS = 100
		cfg.ReqTimeoutMS = 100
		err := cfg.Valid()
		if err != nil {
			panic(err)
		}

		_, err = s.CreateGroup("test", cfg)
		if err != nil {
			logger.Warn("test", "error - %s", err)
		}
	}
}

func ch4tester(id nodeid.NodeId) espnet.Channel {
	r := new(espnet.VChannel)
	r.InitVChannel(fmt.Sprintf("test_%d", id))
	r.Sender = func(msg *espnet.Message) error {
		logger.Debug("test", "sendto %d - %s", id, msg.Dump())
		return nil
	}
	return r
}

func TestUsecase(t *testing.T) {
	cfile := "../../../../bin/config/cluster1-config.json"

	sqliteServer := sqlite.NewSqliteServer("sqliteServer")
	sqliteServer.DefaultBoot()

	nodeId := nodeid.NewService("espnode", sqliteServer)
	boot.QuickDefine(nodeId, "", true)

	ngService := NewService("nodeGroupService", sqliteServer, nodeId)
	boot.QuickDefine(ngService, "", true)

	f1 := func() {
		nid := nodeid.NodeId(2)
		ng := ngService.GetGroup("test")
		cs := new(election.CandidateState)
		cs.Id = nid
		cs.Epoch = 0
		cs.Leader = 0
		cs.Status = election.STATUS_IDLE

		ng.Join(ch4tester(nid), cs)
	}
	if f1 != nil {
	}

	funlist := []func(){
		func4tester(ngService),
		f1,
	}

	boot.TestGo(cfile, 3, funlist)
}
