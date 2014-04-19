package tankbat

type Player struct {
	sch     *ServiceChannel
	teamId  int
	teamNum int
	tankId  int
}

func (this *Player) Id() uint32 {
	return this.sch.Id()
}
