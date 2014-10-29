package acl

import (
	"context"
)

type key int

const key4user key = 0

func BindUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, key4user, user)
}

func UserFromContext(ctx context.Context) *User {
	v := ctx.Value(key4user)
	if v != nil {
		r, ok := v.(*User)
		if ok {
			return r
		}
	}
	return nil
}
