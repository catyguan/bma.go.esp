package smmapi4config

import (
	"bmautil/valutil"
	"boot"
	"config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"smmapi"
	"strings"
)

type smmObject int

func (this smmObject) GetInfo() (*smmapi.SMInfo, error) {
	r := new(smmapi.SMInfo)
	r.Title = "Config"

	r.Actions = make([]*smmapi.SMAction, 0)
	if true {
		a := new(smmapi.SMAction)
		a.Id = "config.detail"
		a.Title = "Detail"
		a.Type = smmapi.SMA_HTTPUI
		a.UIN = "go.config/smm.ui:detail.gl.lua"
		r.Actions = append(r.Actions, a)
	}
	return r, nil
}

func (this smmObject) ExecuteAction(aid string, param map[string]interface{}) (interface{}, error) {
	switch aid {
	case "config.list":
		dir, fname := filepath.Split(boot.StartConfigFile)
		dfile, err0 := os.Open(dir)
		if err0 != nil {
			return nil, err0
		}
		defer dfile.Close()
		filist, err1 := dfile.Readdir(-1)
		if err1 != nil {
			return nil, err1
		}
		r := make([]map[string]interface{}, 0)
		for _, fi := range filist {
			if fi.IsDir() {
				continue
			}
			n := fi.Name()
			if strings.HasSuffix(n, ".json") {
				o := make(map[string]interface{})
				o["name"] = n
				o["main"] = fname == n
				r = append(r, o)
			}
		}
		return r, nil
	case "config.new":
		n := valutil.ToString(param["name"], "")
		if n == "" {
			return nil, fmt.Errorf("miss param 'name'")
		}
		n = n + ".json"
		dir := filepath.Dir(boot.StartConfigFile)
		fname := filepath.Join(dir, n)
		fi, _ := os.Stat(fname)
		if fi != nil {
			return nil, fmt.Errorf("config file '%s' exists", n)
		}
		file, err1 := os.Create(fname)
		if err1 != nil {
			return nil, err1
		}
		defer file.Close()
		_, err2 := file.WriteString("")
		if err2 != nil {
			return nil, err2
		}
		return true, nil
	case "config.update":
		n := valutil.ToString(param["name"], "")
		content := valutil.ToString(param["content"], "")
		if n == "" {
			return nil, fmt.Errorf("miss param 'name'")
		}
		n = n + ".json"
		dir := filepath.Dir(boot.StartConfigFile)
		fname := filepath.Join(dir, n)
		_, err0 := os.Stat(fname)
		if err0 != nil {
			return nil, fmt.Errorf("stat file '%s' fail - %s", n, err0)
		}
		err1 := ioutil.WriteFile(fname, []byte(content), os.ModePerm)
		if err1 != nil {
			return nil, err1
		}
		return true, nil
	case "config.view":
		n := valutil.ToString(param["name"], "")
		if n == "" {
			return nil, fmt.Errorf("miss param 'name'")
		}
		n = n + ".json"
		dir := filepath.Dir(boot.StartConfigFile)
		fname := filepath.Join(dir, n)
		content, err1 := ioutil.ReadFile(fname)
		if err1 != nil {
			return nil, err1
		}
		return string(content), nil
	case "config.parse":
		co, err := config.InitConfig(boot.StartConfigFile)
		if err != nil {
			return fmt.Sprintf("parse fail - %s", err.Error()), nil
		}
		bs, err1 := json.Marshal(co)
		if err1 != nil {
			return nil, err1
		}
		return string(bs), nil
	}
	return nil, fmt.Errorf("unknow action(%s)", aid)
}

func init() {
	smmapi.Add("go.config", smmObject(0))
}
