package context

import "sync/atomic"

type key int

const key4execId key = 0

var (
	execId uint32
)

func CreateExecId(ctx Context) (Context, uint32) {
	id := atomic.AddUint32(&execId, 1)
	if id == 0 {
		id = atomic.AddUint32(&execId, 1)
	}
	r := WithValue(ctx, key4execId, id)
	return r, id
}

func ExecIdFromContext(ctx Context) (uint32, bool) {
	v := ctx.Value(key4execId)
	if v != nil {
		r, ok := v.(uint32)
		return r, ok
	}
	return 0, false
}
