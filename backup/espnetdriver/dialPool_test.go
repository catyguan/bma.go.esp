package espnetdriver

import (
	"boot"
	"fmt"
	"testing"
	"time"
)

func TestDialPoolProperties(t *testing.T) {
	cfg := make(map[string]string)
	cfg["max"] = "5"
	cfg["init"] = "1"
	cfg["coder"] = "espnet"
	cfg["address"] = "127.0.0.1:1080"
	cfg["retry.inc"] = "128"

	p := new(DialPoolSource)
	props := p.GetProperties()
	for _, prop := range props {
		if v, ok := cfg[prop.Name]; ok {
			prop.Value.Commit(v)
		}
	}
	m := p.ToMap()
	fmt.Println(m, p.config)
	fmt.Println(p.ToURI())

	if m == nil {
		p2 := new(DialPoolSource)
		p2.FromMap(m)
		fmt.Println(p2.config)
		fmt.Println(p2.config.Retry)
	}

	fac, err := p.CreateChannelFactory(nil, "test", true)
	fmt.Println(fac, err)
	time.Sleep(time.Duration(5) * time.Second)

	fmt.Printf("before close - %s\n", fac)
	boot.RuntimeStopCloseClean(fac, true)
	time.Sleep(time.Duration(1) * time.Millisecond)
	fmt.Printf("after  close - %s\n", fac)
}
