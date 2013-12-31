package bombman

import "fmt"

type Player struct {
	id      int
	name    string
	channel *ServiceChanel

	MapPos
	died bool

	action     int
	actionDir  int
	actionTime int64
}

func (this *Player) View() string {
	return fmt.Sprintf("%d", this.id)
}

func (this *Player) String() string {
	return fmt.Sprintf("PL%d", this.id)
}

type playerList []*Player

func (ms playerList) Len() int {
	return len(ms)
}

func (ms playerList) Less(i, j int) bool {
	return ms[i].actionTime < ms[j].actionTime
}

func (ms playerList) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}
