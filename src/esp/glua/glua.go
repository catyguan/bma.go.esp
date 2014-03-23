package glua

import (
	"bmautil/goo"
	"boot"
	"lua51"
)

const (
	tag = "glua"
)

type ConfigInfo struct {
}

func (this *ConfigInfo) Valid() error {
	return nil
}

func (this *ConfigInfo) Compare(old *ConfigInfo) int {
	if old == nil {
		return boot.CCR_NEED_START
	}
	return boot.CCR_NONE
}

type GLua struct {
	name   string
	config *ConfigInfo

	goo goo.Goo
	l   *lua51.State
}

func NewGLua(n string, queueSize int, cfg *ConfigInfo) *GLua {
	r := new(GLua)
	r.name = n
	r.config = cfg
	r.goo.InitGoo(tag, queueSize, r.exitHandler)
	return r
}

func (this *GLua) String() string {
	return this.name
}

func (this *GLua) exitHandler() {
	if this.l != nil {
		this.l.Close()
	}
}

func (this *GLua) Run() error {
	if this.goo.GetState() == goo.STATE_INIT {
		this.goo.Run()
	}
	return this.goo.DoSync(func() error {
		return this.doInitLua(cfg)
	})
}

func (this *GLua) doInitLua() error {
	this.l = lua51.NewState()
	return nil
}
