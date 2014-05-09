package mcserver

import (
	"bmautil/valutil"
	"logger"
	"strings"
)

type MemcacheCommand struct {
	Action string
	Params []string
	Data   []byte
}

const (
	MC_RES_NONE     = 0
	MC_RES_ERROR    = 1
	MC_RES_RESPONSE = 2
)

type MemcacheResult struct {
	Response string
	Params   []string
	Data     []byte
}

func (this *MemcacheResult) ToError() (bool, string) {
	if this.Response == "ERROR" || this.Response == "CLIENT_ERROR" || this.Response == "SERVER_ERROR" {
		msg := this.Response
		if len(this.Params) > 2 {
			msg = this.Params[1]
		}
		return true, msg
	}
	return false, ""
}

type MemcacheCoder struct {
	data []byte
	wpos int
}

func NewMemcacheCoder() *MemcacheCoder {
	this := new(MemcacheCoder)
	this.data = make([]byte, 1024)
	return this
}

func (this *MemcacheCoder) Write(data []byte) (n int, err error) {
	l := len(data)
	this.grow(l)
	copy(this.data[this.wpos:], data)
	this.wpos = this.wpos + l
	return l, nil
}

func (this *MemcacheCoder) WriteByte(c byte) error {
	this.grow(1)
	this.data[this.wpos] = c
	this.wpos = this.wpos + 1
	return nil
}

func (this *MemcacheCoder) grow(n int) {
	l := cap(this.data)
	if this.wpos+n > l {
		for {
			l = l + 1024
			if this.wpos+n <= l {
				break
			}
		}
		buf := make([]byte, l)
		copy(buf, this.data)
		this.data = buf
	}
}

func (this *MemcacheCoder) decodeLine() (bool, string, int) {
	var line string
	for i := 0; i < this.wpos; i++ {
		if this.data[i] == '\r' {
			if i+1 < this.wpos {
				if this.data[i+1] == '\n' {
					// find it
					line = string(this.data[0:i])
					return true, line, i + 2
				}
			} else {
				return false, line, 0
			}
		}
	}
	return false, line, 0
}

func (this *MemcacheCoder) DecodeCommand() (bool, *MemcacheCommand) {
	ok, line, rpos := this.decodeLine()
	if !ok {
		return false, nil
	}
	str := strings.TrimSpace(line)
	strlist := strings.Split(str, " ")
	var cmd *MemcacheCommand
	var w string
	if len(strlist) != 0 {
		w = strings.ToLower(strlist[0])

		var data []byte
		switch w {
		case "set", "add", "replace":
			if len(strlist) != 5 {
				return false, nil
			}
			sz := valutil.ToInt(strlist[4], 0)
			if rpos+sz+2 > this.wpos {
				return false, nil
			}

			data = make([]byte, sz)
			copy(data, this.data[rpos:rpos+sz])
			rpos = rpos + sz + 2
		}
		logger.Debug(tag, "memcache command << %s", str)
		cmd = new(MemcacheCommand)
		cmd.Action = w
		cmd.Params = strlist[1:]
		cmd.Data = data

		copy(this.data, this.data[rpos:])
		this.wpos = this.wpos - rpos
	}

	return true, cmd
}

func (this *MemcacheCoder) DecodeResult() (bool, *MemcacheResult) {
	ok, line, rpos := this.decodeLine()
	if !ok {
		return false, nil
	}
	str := strings.TrimSpace(line)

	strlist := strings.Split(str, " ")
	var res *MemcacheResult
	var w string
	if len(strlist) != 0 {
		w = strings.ToUpper(strlist[0])

		var data []byte
		switch w {
		case "VALUE":
			sz := 0
			if len(strlist) < 4 {
				sz = 0
			} else {
				sz = valutil.ToInt(strlist[3], 0)
			}
			if rpos+sz > this.wpos {
				return false, nil
			}

			data = make([]byte, sz)
			copy(data, this.data[rpos:rpos+sz])
			rpos = rpos + sz + 2
		}
		logger.Debug(tag, "memcache result << %s", str)
		res = new(MemcacheResult)
		res.Response = w
		res.Params = strlist[1:]
		res.Data = data

		copy(this.data, this.data[rpos:])
		this.wpos = this.wpos - rpos
	}

	return true, res
}
