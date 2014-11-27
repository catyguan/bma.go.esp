package vmmhttp

import (
	"context"
	"golua"
	"net/http"
)

type key int

const key4req key = 0
const key4resp key = 1

func CreateServ(ctx context.Context, w http.ResponseWriter, req *http.Request) context.Context {
	ctx = context.WithValue(ctx, key4req, req)
	ctx = context.WithValue(ctx, key4resp, w)
	return ctx
}

func RequestFromVM(vm *golua.VM) *http.Request {
	ctx := vm.API_getContext()
	if ctx == nil {
		return nil
	}
	r, _ := RequestFromContext(ctx)
	return r
}

func RequestFromContext(ctx context.Context) (*http.Request, bool) {
	v := ctx.Value(key4req)
	if v != nil {
		r, ok := v.(*http.Request)
		return r, ok
	}
	return nil, false
}

func ResponseFromVM(vm *golua.VM) http.ResponseWriter {
	ctx := vm.API_getContext()
	if ctx == nil {
		return nil
	}
	r, _ := ResponseFromContext(ctx)
	return r
}

func ResponseFromContext(ctx context.Context) (http.ResponseWriter, bool) {
	v := ctx.Value(key4resp)
	if v != nil {
		r, ok := v.(http.ResponseWriter)
		return r, ok
	}
	return nil, false
}
