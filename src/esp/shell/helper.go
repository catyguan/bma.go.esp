package shell

import (
	"fmt"
	"math/rand"
)

type confirmInfo struct {
	word    string
	command string
	param   string
}

func CheckConfirmWithAdminWord(s *Session, cmd, param, word string, adminWord string) bool {
	if adminWord != "" && word == adminWord {
		return true
	}
	return CheckConfirm(s, cmd, param, word)
}

func CheckConfirm(s *Session, cmd, param, word string) bool {
	r := s.Get("@confirmInfo", nil)
	if r != nil {
		info := r.(*confirmInfo)
		return info.command == cmd && info.param == param && info.word == word
	}
	return false
}

func CreateConfirm(s *Session, cmd, p string) string {
	r := s.Get("@confirmInfo", func() interface{} {
		return new(confirmInfo)
	})

	word := fmt.Sprintf("%d", rand.Intn(999999-100000)+100000)

	info := r.(*confirmInfo)
	info.command = cmd
	info.param = p
	info.word = word
	return word
}
