package espnet

type ChannelListener func(ch Channel) error

// 通道的业务界面
type Channel interface {
	Id() uint32

	Name() string

	String() string

	// 关闭
	AskClose()

	// 获取属性/设置属性
	GetProperty(name string) (interface{}, bool)
	SetProperty(name string, val interface{}) bool

	// 上行的接收器
	SetMessageListner(rec MessageListener)

	SendMessage(ev *Message) error

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
	IsBreak() *bool
}
