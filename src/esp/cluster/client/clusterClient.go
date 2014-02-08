package client

import (
	"esp/cluster/clusterbase"
	"esp/espnet"
)

type ClusterClient struct {
	cc    *espnet.ChannelClient
	coder clusterbase.OpCoder
}

func NewClusterClient(cc *espnet.ChannelClient, coder clusterbase.OpCoder, h clusterbase.OpHandler) *ClusterClient {
	r := new(ClusterClient)
	r.cc = cc
	r.coder = coder
	r.handler = h
	return r
}

func (this *ClusterClient) Execute(op interface{}) error {

	return nil
}
