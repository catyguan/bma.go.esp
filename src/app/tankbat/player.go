package tankbat

type Player struct {
	sch *ServiceChannel
}

func (this *Player) Id() uint32 {
	return this.sch.Id()
}
