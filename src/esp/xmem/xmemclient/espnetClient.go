package xmem

import "esp/espnet"

type Client struct {
	c *espnet.ChannelClient
}

func NewClient(c *espnet.ChannelClient) Client {
	this := new(Client)
	this.c = c
	return this
}
