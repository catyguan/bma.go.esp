package authclient

import (
	"bmautil/socket"
	"esp/espnet"
	"fmt"
	"testing"
	"time"
)

func TestUseCase(t *testing.T) {
	cfg := new(espnet.DialConfig)
	cfg.Address = "127.0.0.1:1080"
	pool := espnet.NewDialPool("pool", cfg, 1, func(s *socket.Socket) error {
		s.Trace = 128
		return nil
	})
	pool.Start()
	pool.Run()

	ch, err := pool.NewChannelFactory(espnet.SOCKET_CHANNEL_CODER_ESPNET, 1*time.Second).NewChannel()
	if err != nil {
		t.Error(err)
	} else {
		defer ch.AskClose()
		auth := NewAuthClient(ch, espnet.NewAddress("auth"))

		authToken, err := auth.Login("test", "test")
		if err != nil {
			t.Error(err)
		} else {
			// use authToken
			msg := espnet.NewMessage()
			authToken.Bind(msg)
			fmt.Println(msg.Dump())
			// reply message
			// responseMsg := espnet.NewMessage()
			// authErr := authToken.Error(responseMsg)
			// if authErr != nil {
			// 	authToken := auth.Login("username", "password")
			// }
		}
	}

	pool.Close()
	time.Sleep(100 * time.Millisecond)
}
