package aclserv

import (
	"acl"
	"boot"
	"fmt"
	"strings"
	"testing"
)

func TestServiceBoot(t *testing.T) {
	s := NewService("test")
	boot.Add(s, "", false)

	f := func() {
		s := acl.DumpRuleTree()
		fmt.Println(s)

		acl.DebugMode = true

		user := acl.NewUserS("user@127.0.0.2")
		ps := strings.Split("/a/b/c/d", "/")
		fmt.Println("paths = ", ps)
		ok, rule, err := acl.Check(user, ps, nil, false)
		fmt.Printf("user = %s,rule=%v, result=%v, err = %v\n", user, rule, ok, err)
	}
	boot.TestGo("service_test.json", 3, []func(){f})
}
