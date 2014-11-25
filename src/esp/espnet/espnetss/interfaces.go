package espnetss

import "esp/espnet/espsocket"

const (
	tag = "espnetsdk"
)

var (
	gHandlers map[string]LoginHandler
)

func init() {
	gHandlers = make(map[string]LoginHandler)
	RegisterLoginHandler("none", noneLogin)
}

func RegisterLoginHandler(typ string, lg LoginHandler) {
	gHandlers[typ] = lg
}

func GetLoginHandler(typ string) LoginHandler {
	if lh, ok := gHandlers[typ]; ok {
		return lh
	}
	return nil
}

type LoginHandler func(sock *espsocket.Socket, user string, cert string) (bool, error)
