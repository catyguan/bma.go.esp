package protocol

import (
	"fmt"
	"net"
	"os"
	"testing"
	"thrift"
	"time"
)

func TestClient(t *testing.T) {
	startTime := currentTimeMillis()
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory(), nil)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocket(net.JoinHostPort("127.0.0.1", "9090"))
	transport.SetTimeout(time.Duration(5) * time.Second)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address:", err)
		t.Error("error1")
		return
	}

	useTransport := transportFactory.GetTransport(transport)
	client := NewRpcServiceClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to 127.0.0.1:19090", " ", err)
		t.Error("error2")
		return
	}
	defer transport.Close()

	for i := 0; i < 3; i++ {
		paramMap := make(map[string]string)
		paramMap["name"] = "qinerg"
		paramMap["passwd"] = "123456"
		r1, e1 := client.FunCall(currentTimeMillis(), "login", paramMap)
		fmt.Println(i, "Call->", r1, e1)
	}

	endTime := currentTimeMillis()
	fmt.Println("Program exit. time->", endTime, startTime, (endTime - startTime))
}

// 转换成毫秒
func currentTimeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}
