package n2n

import (
	"bmautil/socket"
	"esp/cluster/nodeinfo"
	"esp/espnet/esnp"
	"esp/espnet/espchannel"
	"logger"
	"sync"
	"time"
)

const (
	tag = "n2n"
)

type remoteInfo struct {
	id      nodeinfo.NodeId
	name    string
	url     string
	channel espchannel.Channel
	pool    *socket.DialPool
	passive bool
}

type Service struct {
	name   string
	config *configInfo

	lock    sync.RWMutex
	remotes []*remoteInfo
	i2r     map[nodeinfo.NodeId]*remoteInfo
	n2r     map[string]*remoteInfo
}

func NewService(n string) *Service {
	this := new(Service)
	this.name = n
	this.remotes = make([]*remoteInfo, 0)
	this.i2r = make(map[nodeinfo.NodeId]*remoteInfo)
	this.n2r = make(map[string]*remoteInfo)
	return this
}

func (this *Service) isLive(ri *remoteInfo) bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	for _, o := range this.remotes {
		if o == ri {
			return true
		}
	}
	return false
}

func (this *Service) onChannelClose(ri *remoteInfo) {
	if !this.isLive(ri) {
		return
	}
	logger.Debug(tag, "'%s' break,try reconnect", ri.name)
	ri.channel = nil
	go func() {
		this.reconnect(ri)
	}()
}

func (this *Service) reconnect(ri *remoteInfo) {
	i := 0
	for {
		if !this.isLive(ri) {
			logger.Debug(tag, "'%s' not live, exit reconnect", ri.name)
			return
		}
		if i%30 == 0 {
			logger.Debug(tag, "'%s' try GetSocket", ri.name)
		}
		i++
		sock, err := ri.pool.GetSocket(1*time.Second, false)
		if err == nil {
			logger.Debug(tag, "'%s' connected", ri.name)
			ri.channel = espchannel.NewSocketChannel(sock, espchannel.SOCKET_CHANNEL_CODER_ESPNET)
			ri.channel.SetCloseListener(this.name, func() {
				this.onChannelClose(ri)
			})
			this.onChannelConnect(ri)
			return
		}
	}
}

func (this *Service) checkAndConnect(n string, url string) {
	this.lock.RLock()
	for _, ri := range this.remotes {
		if ri.name == n {
			this.lock.RUnlock()
			logger.Debug(tag, "'%s' exists, skip connect", n)
			return
		}
	}
	this.lock.RUnlock()

	logger.Debug(tag, "'%s' begin connecting", n)

	ri := new(remoteInfo)
	ri.name = n
	ri.url = url
	this.remotes = append(this.remotes, ri)

	go func() {
		addr, _ := esnp.ParseAddress(ri.url)
		cfg := new(socket.DialPoolConfig)
		cfg.Dial.Address = addr.GetHost()
		cfg.MaxSize = 1
		cfg.InitSize = 1
		pool := socket.NewDialPool(this.name+"_"+n+"_pool", cfg, nil)
		ri.pool = pool
		pool.Start()
		pool.Run()
		this.reconnect(ri)
	}()

	return
}

func (this *Service) _removeAndClose(ri *remoteInfo) {
	if ri.pool != nil {
		ri.pool.AskClose()
		ri.pool = nil
	}
	if ri.channel != nil {
		ri.channel.SetCloseListener(this.name, nil)
		ri.channel.AskClose()
		ri.channel = nil
	}
	if ri.id > 0 {
		delete(this.i2r, ri.id)
	}
	if ri.name != "" {
		delete(this.n2r, ri.name)
	}
	for i, item := range this.remotes {
		if item == ri {
			c := len(this.remotes)
			if c > 1 {
				this.remotes[i], this.remotes[c-1] = this.remotes[i-1], nil
				this.remotes = this.remotes[0 : c-1]
			} else {
				this.remotes = make([]*remoteInfo, 0)
			}
			return
		}
	}
}

func (this *Service) closeRemote(n string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for _, ri := range this.remotes {
		if ri.name == n {
			this._removeAndClose(ri)
			return
		}
	}
}

func (this *Service) closeAllRemote() {
	this.lock.Lock()
	defer this.lock.Unlock()
	tmp := this.remotes
	this.remotes = make([]*remoteInfo, 0)
	for _, ri := range tmp {
		this._removeAndClose(ri)
	}
}

func (this *Service) onChannelConnect(ri *remoteInfo) {
	// send req

}

func (this *Service) doJoin(req *joinReq, ch espchannel.Channel) error {
	return nil
}
