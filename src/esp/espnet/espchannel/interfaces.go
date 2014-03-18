package espchannel

import "esp/espnet/esnp"

type SupportProp interface {
	GetProperty(name string) (interface{}, bool)
	SetProperty(name string, val interface{}) bool
}

type ChannelListener func(ch Channel) error
type ChannelSendCallback func(err error)

// 通道的业务界面
type Channel interface {
	Id() uint32

	Name() string

	String() string

	// 关闭
	AskClose()
	ForceClose()

	// 获取属性/设置属性
	GetProperty(name string) (interface{}, bool)
	SetProperty(name string, val interface{}) bool

	// 上行的接收器
	SetMessageListner(rec esnp.MessageListener)

	PostMessage(ev *esnp.Message) error
	SendMessage(ev *esnp.Message, cb ChannelSendCallback) error

	SetCloseListener(name string, lis func()) error
}

// ChannelFactory
type ChannelFactory interface {
	NewChannel() (Channel, error)
}

type ChannelAcceptor interface {
	SetChannelListener(lis ChannelListener)
}

// ChannelBreakSupport
type BreakSupport interface {
	IsBreak() bool
}
