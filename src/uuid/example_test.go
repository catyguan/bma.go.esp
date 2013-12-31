package uuid

import (
	"fmt"
	"testing"
)

func TestExampleNewV4(t *testing.T) {
	u4, err := NewV4()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(u4)
}

func TestExampleNewV5(t *testing.T) {
	u5, err := NewV5(NamespaceURL, []byte("nu7hat.ch"))
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(u5)
}

func TestExampleParseHex(t *testing.T) {
	u, err := ParseHex("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(u)
}
