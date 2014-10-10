package golua

type GoFunction interface {
	Exec(vm *VM) (int, error)
	String() string
}

type supportFuncName interface {
	FuncName() string
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
	Get() (interface{}, error)
	Set(v interface{}) (bool, error)
}
