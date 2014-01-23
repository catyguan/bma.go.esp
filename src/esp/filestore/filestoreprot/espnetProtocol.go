package filestoreprot

import (
	"crypto/md5"
	"fmt"
	"io"
)

type SessionToken string

type AppId string

type FileStore interface {
	BeginSession(appId AppId) (SessionToken, error)

	SendFile(token SessionToken, path string, content []byte, vcode string) error

	RemoveFile(token SessionToken, path string, vcode string) error

	CommitSession(token SessionToken, vcode string) error

	ReadFile(appId AppId, path string, pos int, readSize int, vcode string) (io.Reader, error)
}

func CreateVCode(token string, key string) string {
	s := token + key
	return fmt.Sprintf("%x", md5.New().Sum([]byte(s)))
}
