package n2n

import (
	"esp/espnet/esnp"
	"fmt"
)

func CheckURL(v string) (*esnp.URL, error) {
	if v == "" {
		return nil, fmt.Errorf("URL invalid")
	}
	o, err := esnp.ParseURL(v)
	if err != nil {
		return nil, fmt.Errorf("URL invalid %s", err)
	}
	host := o.GetHost()
	if host == "" {
		return nil, fmt.Errorf("URL address miss Host")
	}
	return o, nil
}

type RemoteConfigInfo struct {
	Host string
	Code string
}

func (this *RemoteConfigInfo) Valid() error {
	if this.Host == "" {
		return fmt.Errorf("invalid Host")
	}
	return nil
}

func (this *RemoteConfigInfo) Compare(old *RemoteConfigInfo) bool {
	if old == nil {
		return false
	}
	if this.Host != old.Host {
		return false
	}
	if this.Code != old.Code {
		return false
	}
	return true
}

type MapOfRemoteConfigInfo map[string]*RemoteConfigInfo

func (this MapOfRemoteConfigInfo) Valid() error {
	for k, remote := range this {
		err := remote.Valid()
		if err != nil {
			return fmt.Errorf("Remote[%s] %s", k, err)
		}
	}
	return nil
}

func (this MapOfRemoteConfigInfo) Compare(old MapOfRemoteConfigInfo) bool {
	if len(this) != len(old) {
		return false
	}
	for k, ro := range this {
		oro, ok := old[k]
		if ok {
			if !ro.Compare(oro) {
				return false
			}
		}
	}
	for k, _ := range old {
		_, ok := this[k]
		if !ok {
			return false
		}
	}
	return true
}

type ConfigInfo struct {
	Host      string
	Code      string
	Remotes   MapOfRemoteConfigInfo
	TimeoutMS int
}

func (this *ConfigInfo) Valid() error {
	if this.Host == "" {
		return fmt.Errorf("invalid Host")
	}
	if this.TimeoutMS <= 0 {
		this.TimeoutMS = 3000
	}
	if this.Remotes != nil {
		err := this.Remotes.Valid()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *ConfigInfo) Compare(old *ConfigInfo) bool {
	if old == nil {
		return false
	}
	if this.Host != old.Host {
		return false
	}
	if len(this.Remotes) != len(old.Remotes) {
		return false
	}
	if this.Remotes != nil {
		if !this.Remotes.Compare(old.Remotes) {
			return false
		}
	}
	if this.TimeoutMS != old.TimeoutMS {
		return false
	}
	return true
}

func (this *Service) InitConfig(cfg *ConfigInfo) error {
	this.config = cfg
	return nil
}

func (this *Service) Start() bool {
	return this.goo.Run()
}

func (this *Service) Run() bool {
	this.goo.DoSync(func() {
		for k, ro := range this.config.Remotes {
			this.doCheckConnector(k, ro.Host, ro.Code)
		}
	})
	return true
}

func (this *Service) GraceStop(cfg *ConfigInfo) bool {
	for k, oro := range this.config.Remotes {
		var ro *RemoteConfigInfo
		if cfg.Remotes != nil {
			ro = cfg.Remotes[k]
			if !ro.Compare(oro) {
				ro = nil
			}
		}
		if ro != nil {
			continue
		}
		this.goo.DoSync(func() {
			this.doCloseConnector(k)
		})
	}
	return true
}

func (this *Service) Stop() bool {
	this.goo.Stop()
	return true
}

func (this *Service) Cleanup() bool {
	this.goo.StopAndWait()
	return true
}
