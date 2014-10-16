package golua

import "context"

type key int

const key4req key = 0

func CreateRequest(ctx context.Context, req *RequestInfo) context.Context {
	return context.WithValue(ctx, key4req, req)
}

func RequestFromContext(ctx context.Context) (*RequestInfo, bool) {
	v := ctx.Value(key4req)
	if v != nil {
		r, ok := v.(*RequestInfo)
		return r, ok
	}
	return nil, false
}
