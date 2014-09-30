package lua51

import (
	"fmt"
	"testing"
)

func test(L *State) int {
	fmt.Println("hello world! from go!")
	return 0
}

func test2(L *State) int {
	arg := CheckInteger(L, -1)
	argfrombottom := CheckInteger(L, 1)
	fmt.Print("test2 arg: ")
	fmt.Println(arg)
	fmt.Print("from bottom: ")
	fmt.Println(argfrombottom)
	return 0
}

func TestBase(t *testing.T) {
	var L *State

	L = NewState()
	defer L.Close()
	L.OpenLibs()

	L.GetField(LUA_GLOBALSINDEX, "print")
	L.PushString("Hello World!")
	L.Call(1, 0)

	L.PushGoFunction(test)
	L.PushGoFunction(test)
	L.PushGoFunction(test)
	L.PushGoFunction(test)

	L.PushGoFunction(test2)
	L.PushInteger(42)
	L.Call(1, 0)

	L.Call(0, 0)
	L.Call(0, 0)
	L.Call(0, 0)

	// L.Eval("print(1)")
}
