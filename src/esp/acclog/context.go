package acclog

import "context"

type key int

const key4adt key = 0

func CreateAcclogData(ctx context.Context, adt map[string]interface{}) context.Context {
	ctx = context.WithValue(ctx, key4adt, adt)
	return ctx
}

func AcclogDataFromContext(ctx context.Context) (map[string]interface{}, bool) {
	v := ctx.Value(key4adt)
	if v != nil {
		r, ok := v.(map[string]interface{})
		return r, ok
	}
	return nil, false
}
