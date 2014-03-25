package glua

import (
	"fmt"
	"lua51"
	"testing"
)

func TestGoValuesFuncs(t *testing.T) {
	l := lua51.NewState()
	defer l.Close()

	// l.Register("go_int", go_int)
	l.GetGlobal("go_int")
	// l.PushInteger(123)
	l.PushGValue("a123")
	l.PCall(1, 1, 0)
	fmt.Println("here", l.ToString(-1))
}

func TestPushGoValue(t *testing.T) {
	l := lua51.NewState()
	defer l.Close()

	l.PushGValue(123)
	v, id, ok := l.ToGValue(-1)
	l.Pop(1)
	fmt.Println(v, id, ok)
}
