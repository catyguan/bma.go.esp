package main

import (
	"app/tankbat"
	"fmt"
	"time"
)

func main() {
	m := tankbat.NewMatrix(nil, 32)
	m.Run(1, nil, nil)

	for {
		if m.IsClosing() {
			break
		}
		time.Sleep(1 * time.Millisecond)
	}
	time.Sleep(100 * time.Millisecond)
	fmt.Println("END")
}
