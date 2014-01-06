package dchan

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

var (
	startTime int64 = time.Now().UnixNano()
)

func tm() int64 {
	return time.Now().UnixNano() - startTime
}

func TestDChan(t *testing.T) {
	runtime.GOMAXPROCS(5)
	d := NewDChan(1)
	go func() {
		for {
			v, _ := d.Read(nil)
			if v != nil {
				i := v.(int)
				if i < 0 {
					fmt.Println(tm(), "resize", -i)
					d.DoResize(-i)
				}
			}
			fmt.Println(tm(), "R", v)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	for i := 0; i < 2; i++ {
		go func(id int) {
			x := 1
			for {
				fmt.Println(tm(), "W", id, "start")
				x++
				d.Write(x*10 + id)
				fmt.Println(tm(), "W", id, "end")
			}
		}(i)
	}
	time.Sleep(500 * time.Millisecond)
	d.Write(-2)
	time.Sleep(500 * time.Millisecond)

}

func TestDChan2(t *testing.T) {
	runtime.GOMAXPROCS(5)
	d := NewDChan(1)
	go func() {
		for {
			v, _ := d.Read(nil)
			if v != nil {
				i := v.(int)
				if i < 0 {
					fmt.Println(tm(), "resize", -i)
					d.DoResize(-i)
				}
			}
			fmt.Println(tm(), "R", v)
		}
	}()

	for i := 0; i < 10; i++ {
		d.Write(i + 1)
	}
	d.Write(-10)
	for i := 0; i < 10; i++ {
		d.Write(i + 1)
	}
	time.Sleep(500 * time.Millisecond)

}
