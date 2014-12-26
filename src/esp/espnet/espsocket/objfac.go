package espsocket

import (
	"bmautil/conndialpool"
	"bmautil/netutil"
	"bmautil/valutil"
	"fmt"
	"objfac"
	"time"
)

const (
	KIND_SOCKET_PROVIDER      = "socketProvider"
	TYPE_DIAL_SOCKET_PROVIDER = "dial"
	TYPE_POOL_SOCKET_PROVIDER = "pool"
)

func init() {
	objfac.SetObjectFactory(KIND_SOCKET_PROVIDER, TYPE_DIAL_SOCKET_PROVIDER, dialObjectProviderFactory(0))
	objfac.SetObjectFactory(KIND_SOCKET_PROVIDER, TYPE_POOL_SOCKET_PROVIDER, poolSocketProviderFactory(0))
}

////// DialSocketProvider
type DialSocketProvider struct {
	cfg *dialConfig
}

type dialConfig struct {
	Net        string
	Address    string
	MaxPackage int
}

func (this *DialSocketProvider) GetSocket(timeout time.Duration) (Socket, error) {
	conn, err := netutil.DialTimeout(this.cfg.Net, this.cfg.Address, timeout)
	if err != nil {
		return nil, err
	}
	return NewConnSocketN(conn, this.cfg.MaxPackage), nil
}

func (this *DialSocketProvider) Close() {
}

type dialObjectProviderFactory int

func (o dialObjectProviderFactory) Valid(cfg map[string]interface{}, ofp objfac.ObjectFactoryProvider) error {
	var co dialConfig
	if valutil.ToBean(cfg, &co) {
		if co.Net == "" {
			co.Net = "tcp"
		}
		if co.Address == "" {
			return fmt.Errorf("Address empty")
		}
		return nil
	}
	return fmt.Errorf("invalid dialObjectProviderFactory config")
}

func (o dialObjectProviderFactory) Compare(cfg map[string]interface{}, old map[string]interface{}, ofp objfac.ObjectFactoryProvider) (same bool) {
	var co, oo dialConfig
	if !valutil.ToBean(cfg, &co) {
		return false
	}
	if !valutil.ToBean(old, &oo) {
		return false
	}
	if co.Net != oo.Net {
		return false
	}
	if co.Address != oo.Address {
		return false
	}
	if co.MaxPackage != oo.MaxPackage {
		return false
	}
	return true
}

func (o dialObjectProviderFactory) Create(cfg map[string]interface{}, ofp objfac.ObjectFactoryProvider) (interface{}, error) {
	err := o.Valid(cfg, ofp)
	if err != nil {
		return nil, err
	}
	var co dialConfig
	valutil.ToBean(cfg, &co)

	return &DialSocketProvider{&co}, nil
}

////// PoolSocketProvider
type PoolSocketProvider struct {
	cfg  *poolConfig
	pool *conndialpool.DialPool
}

func (this *PoolSocketProvider) GetSocket(timeout time.Duration) (Socket, error) {
	conn, err := this.pool.GetConn(timeout, true)
	if err != nil {
		return nil, err
	}
	return NewConnSocket(conn, this.cfg.MaxPackage), nil
}

func (this *PoolSocketProvider) Close() {
	this.pool.Close()
}

type poolConfig struct {
	Net        string
	Address    string
	MaxPackage int
	TimeoutMS  int
	InitSize   int
	MaxSize    int
	IdleTimeMS int
}

type poolSocketProviderFactory int

func (o poolSocketProviderFactory) Valid(cfg map[string]interface{}, ofp objfac.ObjectFactoryProvider) error {
	var co poolConfig
	if valutil.ToBean(cfg, &co) {
		if co.Net == "" {
			co.Net = "tcp"
		}
		if co.Address == "" {
			return fmt.Errorf("Address empty")
		}
		return nil
	}
	return fmt.Errorf("invalid poolSocketProviderFactory config")
}

func (o poolSocketProviderFactory) Compare(cfg map[string]interface{}, old map[string]interface{}, ofp objfac.ObjectFactoryProvider) (same bool) {
	var co, oo poolConfig
	if !valutil.ToBean(cfg, &co) {
		return false
	}
	if !valutil.ToBean(old, &oo) {
		return false
	}
	if co.Net != oo.Net {
		return false
	}
	if co.Address != oo.Address {
		return false
	}
	if co.TimeoutMS != oo.TimeoutMS {
		return false
	}
	if co.InitSize != oo.InitSize {
		return false
	}
	if co.MaxSize != oo.MaxSize {
		return false
	}
	if co.IdleTimeMS != oo.IdleTimeMS {
		return false
	}
	if co.MaxPackage != oo.MaxPackage {
		return false
	}
	return true
}

func (o poolSocketProviderFactory) Create(cfg map[string]interface{}, ofp objfac.ObjectFactoryProvider) (interface{}, error) {
	err := o.Valid(cfg, ofp)
	if err != nil {
		return nil, err
	}
	var co poolConfig
	valutil.ToBean(cfg, &co)

	pcfg := new(conndialpool.DialPoolConfig)
	pcfg.Net = co.Net
	pcfg.Address = co.Address
	pcfg.InitSize = co.InitSize
	pcfg.MaxSize = co.MaxSize
	if pcfg.MaxSize <= 0 {
		pcfg.MaxSize = 128
	}
	pcfg.IdleMS = co.IdleTimeMS
	pcfg.Valid()
	n := fmt.Sprintf("poolSP_%d", time.Now().Unix())
	pool := conndialpool.NewDialPool(n, pcfg)

	r := new(PoolSocketProvider)
	r.cfg = &co
	r.pool = pool
	return r, nil
}
