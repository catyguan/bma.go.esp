package config

import (
	"os"
	"testing"
)

type block struct {
	Value int
	Nouse bool
}

func TestConfig(t *testing.T) {
	var wd, _ = os.Getwd()
	if wd == "" {
		t.Log("a")
	}
	var err = InitConfig("../../test/test-config.json")
	if err != nil {
		t.Error("init fail", err)
	}
	v1 := GetIntConfig("abc", 1)
	if v1 != 1 {
		t.Error("GetConfig abc !=1,", v1)
	}
	v2 := GetBoolConfig("global.Debug", false)
	if !v2 {
		t.Error("GetConfig debug not true,", v2)
	}
	v3 := GetIntConfig("block.Value", 0)
	if v3 != 123 {
		t.Error("GetConfig block.Value not 123,", v3)
	}
	var v4 block
	if GetBeanConfig("block", &v4) {
		if v4.Value != 123 {
			t.Error("GetConfig block.value not 123,", v4)
		}
	} else {
		t.Error("GetBeanConfig block fail")
	}
	v5 := GetStringConfig("block.Dir", "nil")
	t.Log("block.Dir", v5)

	v6 := GetIntConfig("config2.Id", 0)
	if v6 != 1234 {
		t.Error("include config2 fail, config2.Id not 123,", v6)
	}
}
