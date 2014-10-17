package valutil

import (
	"fmt"
	"testing"
)

func TestValueUtil(t *testing.T) {
	if false {
		v1 := ToBool("true", false)
		if !v1 {
			t.Error("ToBool 'true' => ", v1)
		}
	}
	if true {
		v1 := ToString(3, "")
		fmt.Printf(v1)
	}
}

func T2estSizeUtil(t *testing.T) {
	if true {
		v := 12345678
		s := SizeString(uint64(v), 1024, SizeK)
		if s != "12056.326K" {
			t.Error(s)
		}
	}

	if true {
		s := "200M"
		v, _ := ToSize(s, 1024, SizeB)
		if v != 209715200 {
			t.Errorf("%d", v)
		}
	}

	if true {
		s := "200"
		v, _ := ToSize(s, 1024, SizeB)
		if v != 200 {
			t.Errorf("%d", v)
		}
	}

	if true {
		s := "200M"
		v, _ := ToSize(s, 1024, SizeK)
		if v != 204800 {
			t.Errorf("%d", v)
		}
	}
}

type bean struct {
	DataSource string
	Query      string
	Format     string // default: json, string
}

func T2estBeanToMap(t *testing.T) {
	var o bean
	o.DataSource = "mysql"
	o.Query = "select"
	o.Format = "json"

	m := BeanToMap(o)
	if m == nil {
		t.Error(m)
	} else {
		t.Log(m)
	}

	var o2 bean
	if !ToBean(m, &o2) {
		t.Error(o2)
	} else {
		t.Log(o2)
	}
}
