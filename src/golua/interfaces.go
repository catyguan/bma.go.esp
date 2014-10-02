package golua

type GoFunction interface {
	Exec(vm *VM) (int, error)
	String() string
}

type ACTRES int

var (
	ACTRES_ERROR    = ACTRES(0)
	ACTRES_NEXT     = ACTRES(1)
	ACTRES_BREAK    = ACTRES(2)
	ACTRES_CONTINUE = ACTRES(3)
	ACTRES_RETURN   = ACTRES(4)
)

type Action interface {
	Line() int
	Children() []Action
	Process(vm *VM) (ACTRES, error)
	String() string
}

type VMVar interface {
	Get() (interface{}, error)
	Set(v interface{}) (bool, error)
}
