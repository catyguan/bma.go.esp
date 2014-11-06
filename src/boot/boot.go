package boot

import (
	"config"
	"flag"
	"fmt"
	"logger"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type Phase uint

const (
	PREPARE Phase = iota
	CHECKCONFIG
	INIT
	START
	RUN
	GRACESTOP
	STOP
	CLOSE
	CLEANUP
)

type phaseActionInfo struct {
	name   string
	object BootObject
	flag   interface{}
}

var (
	configFile        string
	currentConfigFile string
	maxp              int
	stopState         uint32
	chan4Stop         chan bool              = make(chan bool, 3)
	phaseActions      []*phaseActionInfo     = make([]*phaseActionInfo, 0)
	objects           map[string]interface{} = make(map[string]interface{})
)

func sname(name string, o interface{}) string {
	if name == "" {
		return fmt.Sprintf("%p", o)
	}
	return name
}

func AddService(o BootObject) {
	Add(o, "", false)
}

func Add(o BootObject, name string, install bool) {
	if name == "" {
		if n, ok := o.(SupportName); ok {
			name = n.Name()
		}
	}
	if name == "" {
		name = fmt.Sprintf("OBJ_%p", o)
	}
	pa := &phaseActionInfo{name, o, nil}
	phaseActions = append(phaseActions, pa)

	if install {
		Install(name, o)
	}
}

func RuntimeStartRun(o BootObject) bool {
	ctx := new(BootContext)
	ctx.IsRestart = false
	ctx.Config = config.Global
	if !o.Start(ctx) {
		return false
	}
	if !o.Run(ctx) {
		return false
	}
	return true
}

func RuntimeGrace(o BootObject) {
	o.GraceStop(nil)
}

func RuntimeStopCloseClean(o BootObject, wait bool) {
	call := func() {
		o.Stop()
		o.Close()
		o.Cleanup()
	}
	if wait {
		call()
	} else {
		go call()
	}
}

type noResponseChecker struct {
	timer *time.Timer
	name  string
	lock  sync.Mutex
}

var (
	nrChecker noResponseChecker
)

func (this *noResponseChecker) noResponseCallback() {
	fmt.Printf("WARN: %s no response, force os.Exit", this.name)
	os.Exit(2)
}

func (this *noResponseChecker) start() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.timer != nil {
		this.timer.Stop()
	}
	this.timer = time.AfterFunc(10*time.Second, this.noResponseCallback)
}

func (this *noResponseChecker) stop() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.timer != nil {
		this.timer.Stop()
		this.timer = nil
	}
}

func (this *noResponseChecker) alive() {
	this.lock.Lock()
	defer this.lock.Unlock()
	if this.timer != nil {
		this.timer.Reset(10 * time.Second)
	}
}

func Alive() {
	nrChecker.alive()
}

func doAction(phase Phase, o BootObject, ctx *BootContext) (r bool) {
	defer func() {
		nrChecker.stop()
		err := recover()
		if err != nil {
			fmt.Println("ERROR: doAction fail%\n%s\n", err, string(debug.Stack()))
			r = false
		}
	}()
	switch phase {
	case STOP, CLOSE, CLEANUP:
		nrChecker.start()
	}
	switch phase {
	case PREPARE:
		o.Prepare()
		r = true
	case CHECKCONFIG:
		r = o.CheckConfig(ctx)
	case INIT:
		r = o.Init(ctx)
	case START:
		r = o.Start(ctx)
	case RUN:
		r = o.Run(ctx)
	case GRACESTOP:
		r = o.GraceStop(ctx)
	case STOP:
		r = o.Stop()
	case CLOSE:
		r = o.Close()
	case CLEANUP:
		r = o.Cleanup()
	}
	return
}

func doActions(phase Phase, ctx *BootContext) (r bool) {
	r = true
	switch phase {
	case GRACESTOP, STOP, CLOSE, CLEANUP:
		nrChecker.start()
		c := len(phaseActions)
		for i := c - 1; i >= 0; i-- {
			ainfo := phaseActions[i]
			ctx.CheckFlag = ainfo.flag
			if phase == GRACESTOP {
				cr := ctx.CheckResult()
				if cr != nil {
					switch cr.Type {
					case CCR_CHANGE:
						fmt.Printf("'%s' change\n", ainfo.name)
					case CCR_NEED_START:
						fmt.Printf("'%s' restart\n", ainfo.name)
					}
				}
			}
			ar := doAction(phase, ainfo.object, ctx)
			if !ar {
				r = false
			}
		}
	default:
		for _, ainfo := range phaseActions {
			nrChecker.name = ainfo.name
			ctx.CheckFlag = ainfo.flag
			ar := doAction(phase, ainfo.object, ctx)
			if !ar {
				r = false
				fmt.Printf("'%s' return false\n", ainfo.name)
				return
			}
			ainfo.flag = ctx.CheckFlag
		}
	}
	return
}

func doPrepare() bool {
	flag.StringVar(&configFile, "config", "", "config file name")
	flag.IntVar(&maxp, "maxp", 0, "GOMAXPROCS")
	fmt.Println("boot preparing")
	ctx := new(BootContext)
	return doActions(PREPARE, ctx)
}

func doInitAndStart(cfg string) (bool, *BootContext) {
	lcok, co := loadConfig(cfg)
	if !lcok {
		return false, nil
	}
	config.Global = co
	ctx := new(BootContext)
	ctx.Config = co

	fmt.Println("boot checking config")
	if !doActions(CHECKCONFIG, ctx) {
		return false, ctx
	}
	fmt.Println("boot initing")
	if doInit(ctx) {
		fmt.Println("boot starting")
		if doActions(START, ctx) {
			return true, ctx
		}
		fmt.Println("ERROR: boot start fail")
	} else {
		fmt.Println("ERROR: boot init fail")
	}
	return false, ctx
}

func Restart() bool {
	if atomic.CompareAndSwapUint32(&stopState, 0, 2) {
		chan4Stop <- true
		return true
	}
	return false
}

func doRestartInitStartRun(ctx *BootContext) bool {
	ctx.IsRestart = true
	fmt.Println("restart initing")
	if !doActions(INIT, ctx) {
		fmt.Println("ERROR: restart init fail")
		return false
	}
	fmt.Println("restart starting")
	if !doActions(START, ctx) {
		fmt.Println("ERROR: restart start fail")
		return false
	}
	fmt.Println("restart running")
	if !doActions(RUN, ctx) {
		fmt.Println("ERROR: restart run fail")
		return false
	}
	return true
}

func doRestart() bool {
	defer func() {
		atomic.CompareAndSwapUint32(&stopState, 2, 0)
	}()
	lcok, co := loadConfig(currentConfigFile)
	if !lcok {
		return true
	}
	ctx := new(BootContext)
	ctx.IsRestart = true
	ctx.Config = co
	fmt.Println("restart checking config")
	if !doActions(CHECKCONFIG, ctx) {
		fmt.Println("ERROR: restart check config fail, skip!")
		return true
	}
	config.Global = co
	fmt.Println("restart gracestoping")
	doActions(GRACESTOP, ctx)
	if doRestartInitStartRun(ctx) {
		LoadTime = time.Now()
		return true
	}
	fmt.Println("ERROR: doRestart fail, stop!")
	return false
}

func loadConfig(cfg string) (bool, config.ConfigObject) {
	if !flag.Parsed() {
		flag.Parse()
	}
	if configFile != "" {
		cfg = configFile
	}
	currentConfigFile = cfg
	co, err := config.InitConfig(cfg)
	if err != nil {
		return false, nil
	}
	return true, co
}

func doInit(ctx *BootContext) bool {
	// GOMAXPROCS
	maxpv := maxp
	if maxpv <= 0 {
		maxpv = config.Global.GetIntConfig("global.GOMAXPROCS", 0)
	}
	if maxpv > 0 {
		fmt.Printf("GOMAXPROCS => %d\n", maxpv)
		runtime.GOMAXPROCS(maxpv)
	}
	DevMode = config.Global.GetBoolConfig("global.DevMode", false)
	Debug = config.Global.GetBoolConfig("global.Debug", false)
	logger.DebugFlag = Debug

	if !ctx.IsRestart {
		lcfg := logger.Config()
		lcfg.InitLogger()
	}

	return doActions(INIT, ctx)
}

func doRunAndWait(ctx *BootContext) {
	fmt.Println("boot running")
	if !doRun(ctx) {
		fmt.Println("ERROR: boot run fail")
		return
	}
	for {
		<-chan4Stop
		st := atomic.LoadUint32(&stopState)
		if st == 1 {
			return
		}
		if st == 2 {
			if !doRestart() {
				return
			}
		}
	}
}

func doRun(ctx *BootContext) bool {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM)
	go func() {
		s := <-sigc
		if s == syscall.SIGHUP {
			fmt.Println("receive SIG =>", s, "RESTART")
			go func() {
				Restart()
			}()
			return
		}
		fmt.Println("receive SIG =>", s, "STOP")
		Shutdown()
	}()
	return doActions(RUN, ctx)
}

func doStopAndClean(ctx *BootContext) {
	fmt.Println("boot stoping")
	doStop(ctx)
	fmt.Println("boot closing")
	doClose(ctx)
	fmt.Println("boot cleanuping")
	doCleanup(ctx)
	fmt.Println("boot end")
	time.Sleep(1 * time.Millisecond)
}

func doStop(ctx *BootContext) {
	doActions(STOP, ctx)
}

func doClose(ctx *BootContext) {
	doActions(CLOSE, ctx)
}

func doCleanup(ctx *BootContext) {
	doActions(CLEANUP, ctx)
	logger.Close()
}

func TestGo(cfgFile string, endWaitSec int, funl []func()) {
	var ctx *BootContext
	defer func() {
		if ctx == nil {
			ctx = new(BootContext)
			ctx.Config = config.Global
		}
		doStopAndClean(ctx)
		UninstallAll()
	}()
	time.AfterFunc(time.Duration(endWaitSec+5)*time.Second, func() {
		fmt.Println("os exit!!!!")
		os.Exit(-1)
	})

	if !doPrepare() {
		fmt.Println("ERROR: boot prepare fail")
		return
	}
	ok := false
	ok, ctx = doInitAndStart(cfgFile)
	if ok {
		fmt.Println("boot running")
		if !doRun(ctx) {
			fmt.Println("ERROR: boot run fail")
			return
		}
		for _, f := range funl {
			f()
		}
		if endWaitSec > 0 {
			time.Sleep(time.Duration(endWaitSec) * time.Second)
		}
	}
}

func Go(cfgFile string) {
	StartTime = time.Now()

	var ctx *BootContext
	defer func() {
		if ctx == nil {
			ctx = new(BootContext)
			ctx.Config = config.Global
		}
		doStopAndClean(ctx)
		UninstallAll()
		os.Exit(1)
	}()

	if !doPrepare() {
		fmt.Println("ERROR: boot prepare fail")
		return
	}
	ok := false
	ok, ctx = doInitAndStart(cfgFile)
	if ok {
		LoadTime = time.Now()
		doRunAndWait(ctx)
	}
}

func Shutdown() {
	atomic.StoreUint32(&stopState, 1)
	chan4Stop <- true
}

func Install(name string, obj interface{}) {
	if _, ok := objects[name]; ok {
		panic("install " + name + " error, exists")
	}
	objects[name] = obj
}

func Uninstall(name string) {
	delete(objects, name)
}

func UninstallAll() {
	objects = make(map[string]interface{})
}

func ObjectFor(name string) interface{} {
	r, ok := objects[name]
	if ok {
		return r
	}
	return nil
}
