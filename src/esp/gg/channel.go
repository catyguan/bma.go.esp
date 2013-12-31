package gg

import (
	"esp/espnet"
)

type GChannel struct {
	name string
}

func (this *GChannel) Name() string {
	return name
}

func (this *GChannel) Join(groupName string) {

}

func (this &GChannel) Leave(groupName string) {

}

func (this *GChannel) SetReceiver(f func(pack *espnet.PPackage)) {

}

func (this *GChannel) Broadcast(pack *espnet.PPackage) {

}

func (this *GChannel) SendTo(who *Address, pack *espnet.PPackage) {

}

func (this *GChannel) Close() {

}
