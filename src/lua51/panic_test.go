package lua51

import (
	"fmt"
	"testing"
)

func TestPanic(t *testing.T) {

	var L *State

	L = NewState()
	defer L.Close()
	L.OpenLibs()

	currentPanicf := L.AtPanic(nil)
	currentPanicf = L.AtPanic(currentPanicf)
	newPanic := func(L1 *State) int {
		fmt.Println("I AM PANICKING!!!")
		return 0
		// return currentPanicf(L1)
	}

	L.AtPanic(newPanic)

	//force a panic
	L.PushNil()
	L.Call(0, 0)
}
