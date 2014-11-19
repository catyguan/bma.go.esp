package boot

type BootWrap struct {
	name         string
	phaseActions map[Phase]interface{}
}

func NewBootWrap(n string) *BootWrap {
	r := new(BootWrap)
	r.name = n
	return r
}

func (this *BootWrap) Name() string {
	return this.name
}

func (this *BootWrap) initObject() {
	if this.phaseActions == nil {
		this.phaseActions = make(map[Phase]interface{})
	}
}

func (this *BootWrap) get(p Phase) (interface{}, bool) {
	if this.phaseActions == nil {
		return nil, false
	}
	v, b := this.phaseActions[p]
	return v, b
}

func (this *BootWrap) toInterfaces() BootObject {
	return this
}

func (this *BootWrap) SetPrepare(f func()) {
	this.initObject()
	this.phaseActions[PREPARE] = f
}

func (this *BootWrap) Prepare() {
	f, ok := this.get(PREPARE)
	if ok {
		f.(func())()
	}
}

func (this *BootWrap) SetCheckConfig(f func(ctx *BootContext) bool) {
	this.initObject()
	this.phaseActions[CHECKCONFIG] = f
}

func (this *BootWrap) CheckConfig(ctx *BootContext) bool {
	f, ok := this.get(CHECKCONFIG)
	if ok {
		return f.(func(ctx *BootContext) bool)(ctx)
	}
	return true
}

func (this *BootWrap) SetInit(f func(ctx *BootContext) bool) {
	this.initObject()
	this.phaseActions[INIT] = f
}

func (this *BootWrap) Init(ctx *BootContext) bool {
	f, ok := this.get(INIT)
	if ok {
		return f.(func(ctx *BootContext) bool)(ctx)
	}
	return true
}

func (this *BootWrap) SetStart(f func(ctx *BootContext) bool) {
	this.initObject()
	this.phaseActions[START] = f
}

func (this *BootWrap) Start(ctx *BootContext) bool {
	f, ok := this.get(START)
	if ok {
		return f.(func(ctx *BootContext) bool)(ctx)
	}
	return true
}

func (this *BootWrap) SetRun(f func(ctx *BootContext) bool) {
	this.initObject()
	this.phaseActions[RUN] = f
}

func (this *BootWrap) Run(ctx *BootContext) bool {
	f, ok := this.get(RUN)
	if ok {
		return f.(func(ctx *BootContext) bool)(ctx)
	}
	return true
}

func (this *BootWrap) SetGraceStop(f func(ctx *BootContext) bool) {
	this.initObject()
	this.phaseActions[GRACESTOP] = f
}

func (this *BootWrap) GraceStop(ctx *BootContext) bool {
	f, ok := this.get(GRACESTOP)
	if ok {
		return f.(func(ctx *BootContext) bool)(ctx)
	}
	return true
}

func (this *BootWrap) SetStop(f func() bool) {
	this.initObject()
	this.phaseActions[STOP] = f
}

func (this *BootWrap) Stop() bool {
	f, ok := this.get(STOP)
	if ok {
		return f.(func() bool)()
	}
	return true
}

func (this *BootWrap) SetClose(f func() bool) {
	this.initObject()
	this.phaseActions[CLOSE] = f
}

func (this *BootWrap) Close() bool {
	f, ok := this.get(CLOSE)
	if ok {
		return f.(func() bool)()
	}
	return true
}

func (this *BootWrap) SetCleanup(f func() bool) {
	this.initObject()
	this.phaseActions[CLEANUP] = f
}

func (this *BootWrap) Cleanup() bool {
	f, ok := this.get(CLEANUP)
	if ok {
		return f.(func() bool)()
	}
	return true
}
