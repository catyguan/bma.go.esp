package smmapi4config

import (
	"boot"
	"fmt"
	"path/filepath"
	"testing"
)

func TestAPI(t *testing.T) {
	boot.StartConfigFile, _ = filepath.Abs("../../../bin/config/glserver-config.json")

	o := smmObject(0)
	param := make(map[string]interface{})
	aid := ""
	if false {
		aid = "config.list"
		rs, err := o.ExecuteAction(aid, param)
		if err != nil {
			t.Error(err)
		} else {
			fmt.Printf("RESULT = %v\n", rs)
		}
		return
	}
	if false {
		aid = "config.new"
		param["name"] = "abc"
		rs, err := o.ExecuteAction(aid, param)
		if err != nil {
			t.Error(err)
		} else {
			fmt.Printf("RESULT = %v\n", rs)
		}
		return
	}
	if true {
		aid = "config.update"
		param["name"] = "abc"
		param["content"] = "test"
		rs, err := o.ExecuteAction(aid, param)
		if err != nil {
			t.Error(err)
		} else {
			fmt.Printf("RESULT = %v\n", rs)
		}
		return
	}
	if false {
		aid = "config.view"
		param["name"] = "glserver-config"
		rs, err := o.ExecuteAction(aid, param)
		if err != nil {
			t.Error(err)
		} else {
			fmt.Printf("RESULT = %v\n", rs)
		}
		return
	}
	if false {
		aid = "config.parse"
		rs, err := o.ExecuteAction(aid, param)
		if err != nil {
			t.Error(err)
		} else {
			fmt.Printf("RESULT = %v\n", rs)
		}
		return
	}
}
