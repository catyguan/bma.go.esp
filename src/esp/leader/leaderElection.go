package leader

type Epoch uint32
type State byte

const (
	STATE_LOOKING   = State(0)
	STATE_LEADING   = State(1)
	STATE_FOLLOWING = State(2)
	STATE_OBSERVING = State(3)
)

type eventInfo struct {
	state  State
	epoch  uint32
	leader uint64
	sender uint64
}

type leaderElection struct {
	service *Service

	epoch  uint32
	leader uint64
}

func newLeaderElection(s *Service) *leaderElection {
	this := new(leaderElection)
	this.service = s
	return this
}

func (this *leaderElection) process(ev eventInfo) {
	switch ev.state {
	case STATE_LOOKING:

	}
}
