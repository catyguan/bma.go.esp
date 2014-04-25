package sgs4rps

import (
	"boot"
	"testing"
)

func TestService(t *testing.T) {
	cfile := "testservice.json"

	s := NewService("RPS")
	boot.Add(s, "", false)

	boot.TestGo(cfile, 5, nil)
}
