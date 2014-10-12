package golua

type GoFunction interface {
	Exec(vm *VM) (int, error)
	IsNative() bool
}

type supportFuncName interface {
	FuncName() (string, string)
}

type ER int

var (
	ER_ERROR    = ER(0)
	ER_NEXT     = ER(1)
	ER_BREAK    = ER(2)
	ER_CONTINUE = ER(3)
	ER_RETURN   = ER(4)
)

type VMVar interface {
	Get(vm *VM) (interface{}, error)
	Set(vm *VM, v interface{}) (bool, error)
}

type VMTable interface {
	Get(vm *VM, key string) (interface{}, error)
	Rawget(key string) interface{}
	Set(key string, val interface{})
	Delete(key string)
	Len() int
	ToMap() map[string]interface{}
}
