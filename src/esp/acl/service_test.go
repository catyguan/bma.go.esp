package acl

import (
	"boot"
	"fmt"
	"testing"
)

func TestServiceBase(t *testing.T) {
	cfg := new(configInfo)

	ulist := make([]*configUserInfo, 0)
	if true {
		u2 := new(configUserInfo)
		u2.Id = "test2"
		u2.Token = "abc"
		u2.Host = []string{"1", "127.0.0.2"}
		ulist = append(ulist, u2)

		u1 := new(configUserInfo)
		u1.Id = "test"
		u1.Token = ""
		u1.Host = []string{"1", "127.0.0.1"}
		ulist = append(ulist, u1)
	}
	cfg.Users = ulist

	olist := make([]*configPriInfo, 0)
	if true {
		p1 := new(configPriInfo)
		p1.Op = "login"
		p1.Who = []string{"*"}
		olist = append(olist, p1)
	}
	cfg.Ops = olist

	cfg.Valid()

	s := NewService("test")
	s.config = cfg

	if true {
		user, err := s.GetUser("test", "", "127.0.0.1")
		fmt.Printf("user = %s, err = %v\n", user, err)
	}

	if true {
		user, err := s.GetUser("test", "", "127.0.0.1")
		fmt.Printf("user = %s, err = %v\n", user, err)
		ok, err2 := s.CheckOp(user, "login")
		fmt.Printf("check op result %v, err = %v\n", ok, err2)
	}
}

func TestServiceBoot(t *testing.T) {
	s := NewService("test")
	boot.Add(s, "", false)

	f := func() {
		user, err := s.GetUser("test", "", "127.0.0.1")
		fmt.Printf("user = %s, err = %v\n", user, err)
	}
	boot.TestGo("service_test.json", 3, []func(){f})
}
