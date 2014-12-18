package esnp

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestMessageId(t *testing.T) {
	p1 := NewMessage()
	MessageLineCoders.MessageId.Set(p1, 88888888)
	b, _ := p1.ToBytes()
	fmt.Println(b)
}

func TestMTXData(t *testing.T) {

	time.AfterFunc(5*time.Second, func() {
		os.Exit(-1)
	})

	p1 := NewMessage()
	MessageLineCoders.XData.Add(p1, 1, 1234, nil)
	MessageLineCoders.XData.Add(p1, 2, "abcdef", nil)
	b, _ := p1.ToBytes()
	fmt.Println(b)

	pr := NewMessageReader()
	pr.Append(b)
	pr.Append(b)
	pr.Append([]byte{1, 2, 3})
	for {
		fmt.Println(pr.buffer[:pr.wpos], pr.rpos, pr.wpos)
		p2, _ := pr.ReadMessage(1024)
		if p2 != nil {
			it := MessageLineCoders.XData.Iterator(p2)
			for ; !it.IsEnd(); it.Next() {
				v, _ := it.Value(nil)
				fmt.Println(it.Xid(), v)
			}
		} else {
			break
		}
	}
}
