package main

import (
	"boot"
	"mcserver"
)

const (
	tag = "mcproxy"
)

func main() {
	cfile := "config/mcproxy-config.json"

	service := NewService("service")
	boot.Add(service, "", false)

	mcp := mcserver.NewMemcacheServer("mcPoint", service)
	boot.Add(mcp, "", false)

	boot.Go(cfile)
}
