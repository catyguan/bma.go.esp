package election

import (
	"esp/cluster/nodeid"
	"logger"
	"math/rand"
	"testing"
	"time"
)

type VM4Test struct {
	s map[nodeid.NodeId]*Superior4Test
}

func (this *VM4Test) Init(sz int) {
	this.s = make(map[nodeid.NodeId]*Superior4Test)
	for i := 0; i < sz; i++ {
		so := new(Superior4Test)
		so.Init(this, i+1)
		this.s[so.id] = so
	}
}

func (this *VM4Test) Run() {
	for _, s := range this.s {
		s.Run()
		func(s *Superior4Test) {
			time.AfterFunc(time.Duration(rand.Intn(20))*time.Millisecond, func() {
				for _, o := range this.s {
					if o != s {
						func(o *Superior4Test) {
							o.DoNow(func() {
								st := new(CandidateState)
								*st = s.c.state
								o.c.JoinPartner(st)
							})
							s.DoNow(func() {
								st := new(CandidateState)
								*st = o.c.state
								s.c.JoinPartner(st)
							})
						}(o)
					}
				}
			})
		}(s)
	}

}

func (this *VM4Test) Stop() {
	for _, s := range this.s {
		s.Stop()
	}
}

func (this *VM4Test) SendTo(id nodeid.NodeId, v interface{}) {
	s, _ := this.s[id]
	time.AfterFunc(time.Duration(rand.Intn(10))*time.Millisecond, func() {
		s.DoNow(v)
	})
}

func (this *VM4Test) SendAll(v interface{}, me nodeid.NodeId) {
	for _, s := range this.s {
		if s.id != me {
			s.DoNow(v)
		}
	}
}

type Superior4Test struct {
	id      nodeid.NodeId
	c       *Candidate
	vm      *VM4Test
	ch      chan interface{}
	invalid bool
}

func (this *Superior4Test) Init(vm *VM4Test, id int) {
	this.vm = vm
	this.id = nodeid.NodeId(id)
	this.c = NewCandidate(this.id, this)
	this.ch = make(chan interface{}, 32)
}

func (this *Superior4Test) Run() {
	go this.doRun()
}

func (this *Superior4Test) Stop() {
	close(this.ch)
}

func (this *Superior4Test) DoNow(f interface{}) {
	if cap(this.ch) > 0 {
		this.ch <- f
	}
}

func (this *Superior4Test) doRun() {
	time.AfterFunc(time.Duration(rand.Intn(10))*time.Millisecond, func() {
		this.DoNow(func() {
			this.c.CheckIdle()
		})
	})
	for {
		time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
		req := <-this.ch
		if req == nil {
			return
		}
		switch v := req.(type) {
		case func(s *Superior4Test):
			v(this)
		case func():
			v()
		case *VoteReq:
			err := this.c.OnVoteReq(v)
			if err != nil {
				logger.Error("test", "process vote fail %s", err)
			}
		case *VoteResp:
			this.c.OnVoteResp(v, nil)
		case *AnnounceReq:
			err := this.c.OnAnnounceReq(v)
			if err != nil {
				logger.Error("test", "process announce fail %s", err)
			}
		case *AnnounceResp:
			this.c.OnAnnounceResp(v, nil)
		}
	}
}

func (this *Superior4Test) Name() string {
	return "Testcase"
}

func (this *Superior4Test) AsyncPostVote(who nodeid.NodeId, vote *VoteReq) {
	this.vm.SendTo(who, vote)
}

func (this *Superior4Test) AsyncRespVote(who nodeid.NodeId, resp *VoteResp) {
	this.vm.SendTo(who, resp)
}

func (this *Superior4Test) AsyncPostAnnounce(who nodeid.NodeId, ann *AnnounceReq) {
	this.vm.SendTo(who, ann)
}

func (this *Superior4Test) AsyncRespAnnounce(who nodeid.NodeId, resp *AnnounceResp) {
	this.vm.SendTo(who, resp)
}

func (this *Superior4Test) DoStartLead(old nodeid.NodeId) error {
	return nil
}

func (this *Superior4Test) DoStartFollow(lid nodeid.NodeId) error {
	return nil
}

func (this *Superior4Test) DoStopFollow() error {
	return nil
}

func (this *Superior4Test) OnCandidateInvalid(id nodeid.NodeId) {

}

func TestBase(t *testing.T) {
	logger.Debug("test", "start")
	rand.Seed(time.Now().UnixNano())

	vm := new(VM4Test)
	vm.Init(5)
	vm.Run()

	time.Sleep(2 * time.Second)
	vm.Stop()
	time.Sleep(1 * time.Millisecond)

	logger.Debug("test", "end")
	for _, s := range vm.s {
		logger.Debug("test", "%s", s.c.state.String())
	}
}
