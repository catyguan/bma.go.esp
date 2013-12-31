package boot

type SupportInit interface {
	Init() bool
}

type SupportStart interface {
	Start() bool
}

type SupportRun interface {
	Run() bool
}

type SupportStop interface {
	Stop() bool
}

type SupportClose interface {
	Close() bool
}

type SupportCleanup interface {
	Cleanup() bool
}

type SupportName interface {
	Name() string
}
