package thriftpoint

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
	"thrift"
)

func open() (thrift.TTransport, *TCacheServerClient, error) {
	host := "localhost"
	portStr := "9090"
	framed := true

	var trans thrift.TTransport
	trans, err := thrift.NewTSocket(net.JoinHostPort(host, portStr))
	if err != nil {
		return nil, nil, err
	}
	if framed {
		trans = thrift.NewTFramedTransport(trans, nil)
	}
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	client := NewTCacheServerClientFactory(trans, protocolFactory)

	if err := trans.Open(); err != nil {
		return nil, nil, err
	}
	return trans, client, nil
}

func TestCacheServerThrift(t *testing.T) {

	trans, client, err := open()
	if err != nil {
		t.Error("open", err)
		return
	}
	defer trans.Close()

	cmd := "load"
	groupName := "c2"
	key := "20729415"
	switch cmd {
	case "get":
		req := new(TCacheRequest)
		req.GroupName = groupName
		req.Key = key
		req.Trace = true
		fmt.Print(client.CacheServerGet(req, nil))
	case "load":
		fmt.Print(client.CacheServerLoad(groupName, key))
		// fmt.Print("\n")
	case "put":
		fmt.Print(client.CacheServerPut(groupName, key, []byte("mytest"), 0))
		// fmt.Print("\n")
	case "delete":
		fmt.Print(client.CacheServerErase(groupName, key))
	}
}

func BenchmarkGet(b *testing.B) {
	trans, client, err := open()
	if err != nil {
		b.Error("open", err)
		return
	}
	defer trans.Close()

	groupName := "c2"
	for i := 0; i < b.N; i++ {
		req := new(TCacheRequest)
		req.GroupName = groupName
		v := rand.Intn(99999-10000) + 10000
		req.Key = fmt.Sprintf("%d", v)
		client.CacheServerGet(req, nil)
	}
}
