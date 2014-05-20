package rtstatus

import (
	"fmt"
	"testing"
)

func TestRTS(t *testing.T) {
	col := NewCollections()
	col.Set("title", "test rts")
	col.Set("name", "hello kitty")
	for i := 0; i < 4; i++ {
		c := col.SubCollections(fmt.Sprintf("address %d", i+1))
		c.Set("city", "mycity")
		c.Set("code", 12345)
	}

	sfm := new(RTSFormatter4String)
	sfm.Init()
	col.Print(sfm)
	fmt.Print(sfm.String())
}
