package sgs4rps

import (
	"esp/espnet/espsocket"
)

type playerAPI interface {
	NewPlayer(sock *espsocket.Socket) (int, error)
	SetNick(psid int, n string) error
	Play(psid int, rps int) error
	Quit(psid int) error
}
