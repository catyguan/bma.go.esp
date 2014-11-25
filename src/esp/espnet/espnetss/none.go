package espnetss

import "esp/espnet/espsocket"

func noneLogin(sock *espsocket.Socket, user string, cert string) (bool, error) {
	return true, nil
}
