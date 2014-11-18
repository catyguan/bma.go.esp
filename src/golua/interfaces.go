package golua

import (
	"errors"
	"fmt"
)

type GoFunction interface {
	Exec(vm *VM, self interface{}) (int, error)
	IsNative() bool
}

type GoObject interface {
	Get(vm *VM, o interface{}, key string) (interface{}, error)
	Set(vm *VM, o interface{}, key string, val interface{}) error
	ToMap(o interface{}) map[string]interface{}
	CanClose() bool
	Close(o interface{})
}

type GoObjectFactory func(vm *VM, n string) (interface{}, error)

type GoLuaInitor func(golua *GoLua)

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

const (
	METATABLE_INDEX    = "__index"
	METATABLE_NEWINDEX = "__newindex"
	KEYWORD_MORE       = "..."
)

type VMTable interface {
	Get(vm *VM, key string) (interface{}, error)
	Rawget(key string) interface{}
	Set(vm *VM, key string, val interface{}) error
	Rawset(key string, val interface{})
	Delete(key string)
	Len() int
	ToMap() map[string]interface{}
}

type VMArray interface {
	Get(vm *VM, idx int) (interface{}, error)
	Set(vm *VM, idx int, val interface{}) error
	Insert(vm *VM, idx int, val interface{}) error
	Add(vm *VM, val interface{}) error
	Delete(vm *VM, idx int) error
	SubArray(start int, end int) ([]interface{}, error)
	Len() int
	ToArray() []interface{}
}

func AssertNil(n string, v interface{}) error {
	if v == nil {
		if n != "" {
			return fmt.Errorf("'%s' null pointer", n)
		}
		return errors.New("null pointer")
	}
	return nil
}

// gofCommon
type GoFunc func(vm *VM, self interface{}) (int, error)
type gofCommon struct {
	name string
	f    GoFunc
}

func NewGOF(n string, f GoFunc) GoFunction {
	r := new(gofCommon)
	r.name = n
	r.f = f
	return r
}

func (this *gofCommon) Exec(vm *VM, self interface{}) (int, error) {
	return this.f(vm, self)
}

func (this *gofCommon) IsNative() bool {
	return true
}

func (this *gofCommon) String() string {
	return fmt.Sprintf("GOF<%s>", this.name)
}

// RequestInfo
type RequestInfo struct {
	Script  string
	Trace   bool
	Context map[string]interface{}
	Data    map[string]interface{}
}

func NewRequestInfo() *RequestInfo {
	r := new(RequestInfo)
	return r
}

func (this *RequestInfo) Valid() error {
	if this.Script == "" {
		return fmt.Errorf("script empty")
	}
	return nil
}

func (this *RequestInfo) String() string {
	return fmt.Sprintf("Req[script=%s,context=%v,data=%v]", this.Script, this.Context, this.Data)
}

const (
	KEY_OBJECT_CONTEXT = "__CONTEXT__"
	KEY_CONTEXT        = "_CONTEXT"
	KEY_REQUEST        = "_REQUEST"
)

type ScriptPreprocess func(content string) (string, error)

type SupportClose interface {
	Close()
}

type SupportTryClose interface {
	TryClose() bool
}
