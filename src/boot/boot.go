package boot

import (
	"config"
	"errors"
	"flag"
	"fmt"
	"logger"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"sort"
	"sync"
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

type PhaseAction func() bool

type phaseActionInfo struct {
	phase  Phase
	name   string
	action PhaseAction
	order  int
}

type phaseActionList []phaseActionInfo

func (ms phaseActionList) Len() int {
	return len(ms)
}

func (ms phaseActionList) Less(i, j int) bool {
	return ms[i].order < ms[j].order
}

func (ms phaseActionList) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

var (
	configFile        string
	currentConfigFile string
	maxp              int
	Chan4Stop         chan bool = make(chan bool, 3)
	callShutdown      sync.Once
	phaseActions      map[Phase]phaseActionList = make(map[Phase]phaseActionList)
	objects           map[string]interface{}    = make(map[string]interface{})
)

func sname(name string, action PhaseAction) string {
	if name == "" {
		return fmt.Sprintf("%p", action)
	}
	return name
}

func Define(ph Phase, name string, action PhaseAction) {
	plist := phaseActions[ph]
	var order int
	if plist == nil {
		order = 0
	} else {
		order = len(plist)
		if ph >= STOP {
			order = -order
		}
	}
	DefineOrder(ph, name, action, order)
}

func find(plist []phaseActionInfo, name string) (int, int) {
	if plist != nil {
		for i, info := range plist {
			if info.name == name {
				return i, info.order
			}
		}
	}
	return -1, 0
}

func DefineAfter(posName string, ph Phase, name string, action PhaseAction) {
	i, order := find(phaseActions[ph], posName)
	if i == -1 {
		panic(errors.New("boot '" + posName + "'' not found"))
	}
	order++
	DefineOrder(ph, name, action, order)
}

func DefineBefore(posName string, ph Phase, name string, action PhaseAction) {
	i, order := find(phaseActions[ph], posName)
	if i == -1 {
		panic(errors.New("boot '" + posName + "'' not found"))
	}
	order--
	DefineOrder(ph, name, action, order)
}

func DefineOrder(ph Phase, name string, action PhaseAction, order int) {
	name = sname(name, action)
	plist := phaseActions[ph]
	info := phaseActionInfo{ph, name, action, order}
	if plist == nil {
		plist = phaseActionList{info}
	} else {
		plist = append(plist, info)
	}
	sort.Sort(plist)
	phaseActions[ph] = plist
}

func QuickDefine(o interface{}, name string, install bool) {
	if name == "" {
		if n, ok := o.(SupportName); ok {
			name = n.Name()
		}
	}
	if name == "" {
		name = fmt.Sprintf("OBJ_%p", o)
	}
	if f, ok := o.(SupportCheckConfig); ok {
		Define(CHECKCONFIG, name, f.CheckConfig)
	}
	if f, ok := o.(SupportInit); ok {
		Define(INIT, name, f.Init)
	}
	if f, ok := o.(SupportStart); ok {
		Define(START, name, f.Start)
	}
	if f, ok := o.(SupportRun); ok {
		Define(RUN, name, f.Run)
	}
	if f, ok := o.(SupportGraceStop); ok {
		Define(GRACESTOP, name, f.GraceStop)
	}
	if f, ok := o.(SupportStop); ok {
		Define(STOP, name, f.Stop)
	}
	if f, ok := o.(SupportClose); ok {
		Define(CLOSE, name, f.Close)
	}
	if f, ok := o.(SupportCleanup); ok {
		Define(CLEANUP, name, f.Cleanup)
	}
	if install {
		Install(name, o)
	}
}

func RuntimeStartRun(o interface{}) bool {
	if f, ok := o.(SupportStart); ok {
		if !f.Start() {
			return false
		}
	}
	if f, ok := o.(SupportRun); ok {
		if !f.Run() {
			return false
		}
	}
	return true
}

func RuntimeGrace(o interface{}) {
	call := func() {
		if f, ok := o.(SupportGraceStop); ok {
			f.GraceStop()
		}
	}
	call()
}

func RuntimeStopCloseClean(o interface{}, wait bool) {
	call := func() {
		if f, ok := o.(SupportStop); ok {
			f.Stop()
		}
		if f, ok := o.(SupportClose); ok {
			f.Close()
		}
		if f, ok := o.(SupportCleanup); ok {
			f.Cleanup()
		}
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

func doAction(phase Phase, action PhaseAction) (r bool) {
	defer func() {
		nrChecker.stop()
		err := recover()
		if err != nil {
			fmt.Sprintf("ERROR: boot fail%\n%s\n", err, debug.Stack())
			r = false
		}
	}()
	switch phase {
	case STOP, CLOSE, CLEANUP:
		nrChecker.start()
	}
	r = action()
	return
}

func doActions(phase Phase, doAll bool) (r bool) {
	r = true
	if alist, ok := phaseActions[phase]; ok {
		for _, ainfo := range alist {
			nrChecker.name = ainfo.name
			if !doAction(phase, ainfo.action) {
				r = false
				if !doAll {
					fmt.Printf("'%s' return false\n", ainfo.name)
					return
				}
			}
		}
	}
	return
}

func Prepare() bool {
	flag.StringVar(&configFile, "config", "", "config file name")
	flag.IntVar(&maxp, "maxp", 0, "GOMAXPROCS")
	fmt.Println("boot preparing")
	return doActions(PREPARE, false)
}

func InitAndStart(cfg string) bool {
	if !loadConfig(cfg) {
		return false
	}
	fmt.Println("boot checking config")
	if !CheckConfig() {
		return false
	}
	fmt.Println("boot initing")
	if Init() {
		fmt.Println("boot starting")
		if Start() {
			return true
		}
		fmt.Println("ERROR: boot start fail")
	} else {
		fmt.Println("ERROR: boot init fail")
	}
	return false
}

func doRestart() bool {
	fmt.Println("restart initing")
	if !doActions(INIT, false) {
		fmt.Println("ERROR: restart init fail")
		return false
	}
	fmt.Println("restart starting")
	if !Start() {
		fmt.Println("ERROR: restart start fail")
		return false
	}
	fmt.Println("restart running")
	if !doActions(RUN, true) {
		fmt.Println("ERROR: restart run fail")
		return false
	}
	return true
}

func Restart() bool {
	tmp := config.ConfigData
	loadConfig(currentConfigFile)
	fmt.Println("restart checking config")
	if !CheckConfig() {
		fmt.Println("ERROR: restart check config fail, skip!")
		config.ConfigData = tmp
		return false
	}
	fmt.Println("restart gracestoping")
	doActions(GRACESTOP, true)
	if !doRestart() {
		fmt.Println("ERROR: doRestart fail, stop!")
		Chan4Stop <- true
	}
	return true
}

func loadConfig(cfg string) bool {
	if !flag.Parsed() {
		flag.Parse()
	}
	if configFile != "" {
		cfg = configFile
	}
	currentConfigFile = cfg
	err := config.InitConfig(cfg)
	if err != nil {
		return false
	}
	return true
}

func CheckConfig() bool {
	return doActions(CHECKCONFIG, false)
}

func Init() bool {
	// GOMAXPROCS
	maxpv := maxp
	if maxpv <= 0 {
		maxpv = config.GetIntConfig("global.GOMAXPROCS", 0)
	}
	if maxpv > 0 {
		fmt.Printf("GOMAXPROCS => %d\n", maxpv)
		runtime.GOMAXPROCS(maxpv)
	}

	lcfg := logger.Config()
	lcfg.InitLogger()

	return doActions(INIT, false)
}

func Start() bool {
	return doActions(START, false)
}

func RunAndWait() {
	fmt.Println("boot running")
	if !Run() {
		fmt.Println("ERROR: boot run fail")
		return
	}
	<-Chan4Stop
}

func Run() bool {
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
		Chan4Stop <- true
	}()
	return doActions(RUN, true)
}

func StopAndClean() {
	fmt.Println("boot stoping")
	Stop()
	fmt.Println("boot closing")
	Close()
	fmt.Println("boot cleanuping")
	Cleanup()
	fmt.Println("boot end")
	time.Sleep(1 * time.Millisecond)
}

func Stop() {
	doActions(STOP, true)
}

func Close() {
	doActions(CLOSE, true)
}

func Cleanup() {
	doActions(CLEANUP, true)
	logger.Close()
}

func TestGo(cfgFile string, endWaitSec int, funl []func()) {
	defer func() {
		StopAndClean()
		UninstallAll()
	}()

	if !Prepare() {
		fmt.Println("ERROR: boot prepare fail")
		return
	}
	if InitAndStart(cfgFile) {
		fmt.Println("boot running")
		if !Run() {
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
	defer func() {
		StopAndClean()
		UninstallAll()
		os.Exit(1)
	}()

	if !Prepare() {
		fmt.Println("ERROR: boot prepare fail")
		return
	}
	if InitAndStart(cfgFile) {
		RunAndWait()
	}
}

func Shutdown() {
	callShutdown.Do(func() {
		Chan4Stop <- true
	})
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
