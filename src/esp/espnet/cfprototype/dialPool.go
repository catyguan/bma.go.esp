package cfprototype

import (
	"bmautil/socket"
	"bmautil/valutil"
	"errors"
	"esp/espnet"
	"fmt"
	"net/url"
	"time"
	"uprop"
)

// DialPoolPrototype
type DialPoolPrototype struct {
	config       *espnet.DialPoolConfig
	channelCoder string
	getTimeout   time.Duration
	trace        int
}

func (this *DialPoolPrototype) Valid() error {
	if this.config == nil {
		return errors.New("config empty")
	}
	if err := this.config.Valid(); err != nil {
		return err
	}
	if this.channelCoder == "" {
		return errors.New("coder empty")
	}
	if espnet.GetSocketChannelCoder(this.channelCoder) == nil {
		return fmt.Errorf("coder(%s) invalid", this.channelCoder)
	}
	return nil
}

func (this *DialPoolPrototype) ToMap() map[string]interface{} {
	if this.config != nil {
		m := valutil.BeanToMap(this.config)
		m["coder"] = this.channelCoder
		m["trace"] = this.trace
		t := time.Duration(this.getTimeout) * time.Millisecond
		m["getTimeout"] = t.Nanoseconds()
		return m
	}
	return nil
}

func (this *DialPoolPrototype) FromMap(data map[string]interface{}) error {
	if data != nil {
		cfg := new(espnet.DialPoolConfig)
		valutil.ToBean(data, cfg)
		this.config = cfg
		this.channelCoder = valutil.ToString(data["coder"], "")
		this.trace = valutil.ToInt(data["trace"], 0)
		t := valutil.ToInt64(data["getTimeout"], 0)
		this.getTimeout = time.Duration(t) * time.Nanosecond

		if err := this.Valid(); err != nil {
			return err
		}
	}
	return nil
}

func (this *DialPoolPrototype) ToURI() (*url.URL, error) {
	if err := this.Valid(); err != nil {
		return nil, err
	}

	v := url.Values{}
	props := this.GetProperties()
	for _, p := range props {
		if p.Name == "address" {
			continue
		}
		if p.Value != nil && p.Value.Value != nil {
			v.Set(p.Name, fmt.Sprintf("%v", p.Value.Value))
		}
	}
	s := fmt.Sprintf("dial://%s/?%s", this.config.Dial.Address, v.Encode())
	return url.Parse(s)
}

func (this *DialPoolPrototype) FromURI(u *url.URL) error {
	props := this.GetProperties()
	q := u.Query()
	for k, _ := range q {
		err := uprop.Helper.Set(props, k, q.Get(k))
		if err != nil {
			return err
		}
	}
	return this.Valid()
}

func (this *DialPoolPrototype) GetProperties() []*uprop.UProperty {
	r := make([]*uprop.UProperty, 0)
	if this.config == nil {
		this.config = new(espnet.DialPoolConfig)
	}
	if this.config.Retry == nil {
		this.config.Retry = this.config.DefaultRetryConfig()
	}
	cfg := this.config
	r = append(r, uprop.NewUProperty("address", cfg.Dial.Address, false, "dial address", func(v string) error {
		cfg.Dial.Address = v
		return nil
	}))
	r = append(r, uprop.NewUProperty("coder", this.channelCoder, false, "socket channel coder", func(v string) error {
		if espnet.GetSocketChannelCoder(v) == nil {
			return fmt.Errorf("coder(%s) invalid", v)
		}
		this.channelCoder = v
		return nil
	}))
	r = append(r, uprop.NewUProperty("max", cfg.MaxSize, false, "pool max size", func(v string) error {
		cfg.MaxSize = valutil.ToInt(v, 0)
		return nil
	}))
	r = append(r, uprop.NewUProperty("net", cfg.Dial.Net, true, "dial net,default tcp", func(v string) error {
		cfg.Dial.Address = v
		return nil
	}))
	r = append(r, uprop.NewUProperty("timeout", cfg.Dial.TimeoutMS, true, "dial timeout in MS, default 5000", func(v string) error {
		cfg.Dial.TimeoutMS = valutil.ToInt(v, 0)
		return nil
	}))
	t2 := int(this.getTimeout.Seconds())
	r = append(r, uprop.NewUProperty("getTimeout", t2, true, "pool get timeout in Sec", func(v string) error {
		this.getTimeout = time.Duration(valutil.ToInt(v, 0)) * time.Second
		return nil
	}))
	r = append(r, uprop.NewUProperty("init", cfg.InitSize, true, "pool keep alive size", func(v string) error {
		cfg.InitSize = valutil.ToInt(v, 0)
		return nil
	}))
	rcfg := this.config.Retry
	r = append(r, uprop.NewUProperty("retry.inc", rcfg.DelayIncrease, true, "retry delay increase in MS, default 200", func(v string) error {
		rcfg.DelayIncrease = valutil.ToInt(v, 0)
		return nil
	}))
	r = append(r, uprop.NewUProperty("retry.start", rcfg.DelayMin, true, "retry start delay in MS, default 100", func(v string) error {
		rcfg.DelayMin = valutil.ToInt(v, 0)
		return nil
	}))
	r = append(r, uprop.NewUProperty("retry.limit", rcfg.DelayMax, true, "retry max delay in MS, default 1000", func(v string) error {
		rcfg.DelayMax = valutil.ToInt(v, 0)
		return nil
	}))
	r = append(r, uprop.NewUProperty("retry.max", rcfg.Max, true, "retry max times, 0 mean no limit", func(v string) error {
		rcfg.Max = valutil.ToInt(v, 0)
		return nil
	}))
	t := int(cfg.RetryFailInfoDruation.Seconds())
	r = append(r, uprop.NewUProperty("retry.info", t, true, "retry fail info duration in Sec, default 30", func(v string) error {
		cfg.RetryFailInfoDruation = time.Duration(valutil.ToInt(v, 0)) * time.Second
		return nil
	}))
	r = append(r, uprop.NewUProperty("sock.trace", this.trace, true, "socket trace", func(v string) error {
		this.trace = valutil.ToInt(v, 0)
		return nil
	}))
	return r
}

func (this *DialPoolPrototype) CreateChannelFactory(storage ChannelFactoryStorage, name string, start bool) (espnet.ChannelFactory, error) {
	if this.config == nil {
		return nil, errors.New("config empty")
	}
	if err := this.Valid(); err != nil {
		return nil, err
	}
	pool := espnet.NewDialPool(name, this.config, this.onSocketAccept)
	if start {
		if !pool.Start() {
			return nil, errors.New("pool start fail")
		}
		if !pool.Run() {
			return nil, errors.New("pool run fail")
		}
	}
	return pool.NewChannelFactory(this.channelCoder, this.getTimeout), nil
}

func (this *DialPoolPrototype) onSocketAccept(sock *socket.Socket) error {
	if this.trace > 0 {
		sock.Trace = this.trace
	}
	return nil
}
