package acclog

import (
	"bytes"
	"fmt"
	"time"
)

func NewFilenameCreator(n string) func(tm time.Time, num int) string {
	return func(tm time.Time, num int) string {
		out := bytes.NewBuffer(make([]byte, 0))
		out.WriteString(n)
		out.WriteString("_")
		out.WriteString(tm.Format("20060102"))
		if num != 0 {
			out.WriteString(".")
			out.WriteString(fmt.Sprintf("%d", num))
		}
		out.WriteString(".log")
		return out.String()
	}
}

type simpleAccLogInfo struct {
	message string
	time    time.Time
}

func (this *simpleAccLogInfo) Message(cfg map[string]string) string {
	return fmt.Sprintf("%s,%s\n", this.time.Format("2006-01-02 15:04:05"), this.message)
}

func (this *simpleAccLogInfo) TimeDay() int {
	return this.time.Day()
}

func NewSimpleLog(msg string) AccLogInfo {
	r := new(simpleAccLogInfo)
	r.message = msg
	r.time = time.Now()
	return r
}

type commonAccLogInfo struct {
	data       map[string]interface{}
	time       time.Time
	timeUseSec float64
}

func (this *commonAccLogInfo) Message(cfg map[string]string) string {
	out := bytes.NewBuffer(make([]byte, 0))
	out.WriteString("t=")
	out.WriteString(this.time.Format("2006-01-02 15:04:05"))
	out.WriteString("`")
	for k, v := range this.data {
		if v != nil {
			out.WriteString(k)
			out.WriteString("=")
			out.WriteString(fmt.Sprintf("%v", v))
			out.WriteString("`")
		}
	}
	out.WriteString("tu=")
	out.WriteString(fmt.Sprintf("%f", this.timeUseSec))
	out.WriteByte('\n')
	return out.String()
}

func (this *commonAccLogInfo) TimeDay() int {
	return this.time.Day()
}

func NewCommonLog(dt map[string]interface{}, tuseSec float64) AccLogInfo {
	r := new(commonAccLogInfo)
	r.data = dt
	r.time = time.Now()
	r.timeUseSec = tuseSec
	return r
}
