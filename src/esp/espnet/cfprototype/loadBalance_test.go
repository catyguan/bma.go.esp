package cfprototype

import (
	"bmautil/valutil"
	"encoding/json"
	"esp/espnet"
	"testing"
	"time"
)

func TestToBean(t *testing.T) {
	cfg := new(LoadBalanceConfig)
	cfg.AddName("f1")
	cfg.AddName("f2")
	cfg.AddName("f3")
	obj := valutil.BeanToMap(cfg)
	t.Error(obj)
	b, _ := json.Marshal(obj)
	t.Error(string(b))
}

func TestLoadBalance(t *testing.T) {
	cfg := new(LoadBalanceConfig)
	cfg.FailSkipTimeMS = 20
	cfg.AddName("f1").FailOver = true
	cfg.AddName("f2").Priority = 2
	cfg.AddName("f3").FailOver = true

	dcfg := new(espnet.DialPoolConfig)
	dcfg.Dial.Address = "127.0.0.1:1080"
	dcfg.MaxSize = 3

	fcfg := new(espnet.DialPoolConfig)
	fcfg.Dial.Address = "127.0.0.1:1081"
	fcfg.MaxSize = 3
	fcfg.InitSize = 1

	df1 := espnet.NewDialPool("df1", dcfg, nil).NewChannelFactory("espnet", 0)
	df2 := espnet.NewDialPool("df2", fcfg, nil).NewChannelFactory("espnet", 0)
	df3 := espnet.NewDialPool("df3", dcfg, nil).NewChannelFactory("espnet", 0)

	sto := new(SimpleChannelFactoryStorage)
	sto.Factory = make(map[string]espnet.ChannelFactory)
	sto.Factory["f1"] = df1
	sto.Factory["f2"] = df2
	sto.Factory["f3"] = df3

	lb := NewLoadBalanceChannelFactory(sto, cfg)

	for i := 0; i < 10; i++ {
		ch, err := lb.NewChannel()
		t.Error(lb, ch, err)
		if ch != nil {
			ch.AskClose()
		}
	}
	time.Sleep(20 * time.Millisecond)
	for i := 0; i < 10; i++ {
		ch, err := lb.NewChannel()
		t.Error(lb, ch, err)
		if ch != nil {
			ch.AskClose()
		}
	}

	time.Sleep(200 * time.Millisecond)
}
