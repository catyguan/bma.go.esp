package goo

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestGoo1(t *testing.T) {
	time.AfterFunc(5*time.Second, func() {
		os.Exit(-100)
	})

	o := new(Goo)
	o.EDebug = true
	o.InitGoo("test1", 16, exithandler4test)
	o.Run()
	o.DoSync(func() {
		fmt.Println("Say hi~~~~")
	})
	o.Do(func() error {
		return fmt.Errorf("error from Goo")
	}, func(err error) {
		fmt.Println("PrintErr", err)
	})
	time.Sleep(1 * time.Second)
	o.StopAndWait()
}

func exithandler4test() {
	fmt.Println("goo exit!!!")
}
