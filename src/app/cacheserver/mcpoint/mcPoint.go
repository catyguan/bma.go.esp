package mcpoint

import (
	"app/cacheserver"
	"bmautil/netutil"
	"bmautil/valutil"
	"boot"
	"config"
	"logger"
	"mcserver"
	"regexp"
	"sync"
	"time"
)

const (
	tag = "mcpoint"
)

type cacheRouter struct {
	pattern string
	matcher *regexp.Regexp
	group   string
}

type MemcachePoint struct {
	name         string
	service      *cacheserver.CacheService
	Version      string
	router       []*cacheRouter
	defaultGroup string
	disable      bool
}

func NewMemcachePoint(name string, s *cacheserver.CacheService) *MemcachePoint {
	this := new(MemcachePoint)
	this.name = name
	this.service = s
	this.Version = "1.0.0"
	this.router = make([]*cacheRouter, 0)
	return this
}

type routerConfig struct {
	Pattern string
	Group   string
}

type configInfo struct {
	Version      string
	Router       []routerConfig
	DefaultGroup string
	Disable      bool
}

func (this *MemcachePoint) Init() bool {
	cfg := configInfo{}
	if config.GetBeanConfig(this.name, &cfg) {
		if cfg.Disable {
			this.disable = cfg.Disable
		} else {
			if cfg.Version == "" {
				this.Version = "1.0.0"
			} else {
				this.Version = this.Version
			}
			if cfg.Router != nil {
				rlist := make([]*cacheRouter, 0, len(cfg.Router))
				for _, r := range cfg.Router {
					cr := new(cacheRouter)
					var err error
					cr.matcher, err = regexp.Compile(r.Pattern)
					if err != nil {
						logger.Error(tag, "compile router regexp '%s' fail - %s", r.Pattern, err)
						return false
					}
					cr.pattern = r.Pattern
					cr.group = r.Group
					rlist = append(rlist, cr)
				}
				this.router = rlist
			}
			this.defaultGroup = cfg.DefaultGroup
		}
	} else {
		this.disable = true
	}
	if this.disable {
		logger.Info(tag, "%s disabled", this.name)
	}
	return true
}

func (this *MemcachePoint) DefaultBoot() {
	boot.Define(boot.INIT, this.name, this.Init)
}

func (this *MemcachePoint) Handle(c *netutil.Channel, cmd *mcserver.MemcacheCommand) (mcserver.HandleCode, error) {
	switch cmd.Action {
	case "version":
		c.Write([]byte("VERSION " + this.Version + "\r\n"))
		return mcserver.DONE, nil
	case "get":
		return this.serveGet(c, cmd)
	case "set":
		return this.servePut(c, cmd)
	case "delete":
		return this.serveErase(c, cmd)
	}
	return mcserver.UNKNOW_COMMAND, nil
}

func (this *MemcachePoint) KeyToGroup(key string) string {
	for _, r := range this.router {
		if r.matcher.MatchString(key) {
			return r.group
		}
	}
	return this.defaultGroup
}

func (this *MemcachePoint) serveGet(c *netutil.Channel, cmd *mcserver.MemcacheCommand) (mcserver.HandleCode, error) {
	if len(cmd.Params) < 1 {
		return mcserver.UNKNOW_COMMAND, nil
	}

	rep := make(chan *cacheserver.GetResult, len(cmd.Params))
	defer close(rep)

	doget := func(key string, rep chan *cacheserver.GetResult) {

		groupName := this.KeyToGroup(key)

		logger.Debug(tag, "serveGet(%v,%v)", groupName, key)

		req := cacheserver.NewGetRequest(key)
		req.TimeoutMs = 0

		err := this.service.Get(groupName, req, rep)
		if err != nil {
			logger.Warn(tag, "CacheServerGet fail - %s", err.Error())
			r := cacheserver.NewGetResult(groupName, key, false)
			r.Fail(err, nil)
			rep <- r
			return
		}
	}

	l := len(cmd.Params)
	if l == 1 {
		doget(cmd.Params[0], rep)
	} else {
		var wg sync.WaitGroup
		for _, key := range cmd.Params {
			wg.Add(1)
			go func(key string, rep chan *cacheserver.GetResult) {
				defer wg.Done()
				doget(key, rep)
			}(key, rep)
		}
		wg.Wait()
	}

	for i := 0; i < l; i++ {
		result := <-rep
		if result == nil {
			continue
		}

		logger.Debug(tag, "CacheServerGet %v -> %v", result.Key, result.Done)
		sz := 0
		if result.Value != nil {
			sz = len(result.Value)
		}
		if result.Done {
			val := result.Value
			flags := uint16(0)
			if sz > 3 {
				if val[0] == 0 {
					flags = uint16(val[1]) | uint16(val[2])<<8
					val = val[3:]
					sz -= 3
				}
			}
			c.Write([]byte(logger.Sprintf("VALUE %s %d %d\r\n", result.Key, flags, sz)))
			if val != nil {
				c.Write(val)
			}
			c.Write([]byte("\r\n"))
		}
	}
	c.Write([]byte("END\r\n"))
	return mcserver.DONE, nil
}

func (this *MemcachePoint) servePut(c *netutil.Channel, cmd *mcserver.MemcacheCommand) (mcserver.HandleCode, error) {
	if len(cmd.Params) < 4 {
		return mcserver.UNKNOW_COMMAND, nil
	}
	key := cmd.Params[0]
	groupName := this.KeyToGroup(key)
	value := cmd.Data
	flags := valutil.ToUint16(cmd.Params[1], 0)
	exptime := valutil.ToInt64(cmd.Params[2], 0)
	if exptime < 0 {
		return mcserver.UNKNOW_COMMAND, nil
	}
	deadtime := int64(0)
	if exptime > 0 {
		if exptime <= 60*60*24*30 {
			deadtime = time.Now().Unix() + exptime
		} else {
			deadtime = exptime
		}
	}
	if flags > 0 {
		val := value
		value = make([]byte, len(val)+3)
		value[0] = 0
		value[1] = byte(flags)
		value[2] = byte(flags >> 8)
		copy(value[3:], val)
		logger.Debug(tag, "merge flags -> %v", value)
	}

	logger.Debug(tag, "servePut(%s,%s,%d,%d)", groupName, key, flags, deadtime)
	err := this.service.Put(groupName, key, value, deadtime)
	if err != nil {
		logger.Warn(tag, "servePut fail - %s", err.Error())
		return mcserver.DONE, err
	}
	c.Write([]byte("STORED\r\n"))
	return mcserver.DONE, nil
}

func (this *MemcachePoint) serveErase(c *netutil.Channel, cmd *mcserver.MemcacheCommand) (mcserver.HandleCode, error) {
	if len(cmd.Params) < 1 {
		return mcserver.UNKNOW_COMMAND, nil
	}
	key := cmd.Params[0]
	groupName := this.KeyToGroup(key)

	logger.Debug(tag, "serveErase(%s,%s)", groupName, key)

	ok, err := this.service.Delete(groupName, key)
	if err != nil {
		logger.Warn(tag, "serveErase fail - %s", err.Error())
		return mcserver.DONE, err
	}
	if ok {
		c.Write([]byte("DELETED\r\n"))
	} else {
		c.Write([]byte("NOT_FOUND\r\n"))
	}
	return mcserver.DONE, nil
}
