package aclmux

import (
	"acl"
	"net/http"
)

type UserProvider func(r *http.Request) (*acl.User, error)
