package filelog

import (
	"bmautil/syncutil"
	"errors"
	"fmt"
	"os"
)

type FileLog struct {
	filename   string
	file       *os.File
	out        chan []byte
	closeState *syncutil.CloseState

	EnablePrint bool
}

func NewFileLog(filename string, chsize int) *FileLog {
	this := new(FileLog)
	this.filename = filename
	this.out = make(chan []byte, chsize)
	this.closeState = syncutil.NewCloseState()
	return this
}

func (this *FileLog) Open() error {
	if this.file != nil {
		return nil
	}
	fd, err := os.OpenFile(this.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		return err
	}
	this.file = fd
	go this.run()
	return nil
}

func (this *FileLog) run() {
	defer func() {
		close(this.out)
		if this.file != nil {
			this.file.Close()
		}
		this.closeState.DoneClose()
	}()
	for {
		b := <-this.out
		if b == nil {
			return
		}
		this.file.Write(b)
		if this.EnablePrint {
			os.Stdout.Write(b)
		}
	}
}

func (this *FileLog) WriteString(s string) error {
	_, err := this.Write([]byte(s))
	return err
}

func (this *FileLog) Printf(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return this.WriteString(msg)
}

func (this *FileLog) Println(s string) error {
	_, err := this.Write([]byte(s))
	if err != nil {
		return err
	}
	_, err = this.Write([]byte{'\n'})
	return err
}

func (this *FileLog) Write(p []byte) (n int, err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = errors.New("closed")
		}
	}()
	if p == nil {
		return 0, nil
	}
	this.out <- p
	return len(p), nil
}

func (this *FileLog) WriteByte(c byte) error {
	_, err := this.Write([]byte{c})
	return err
}

func (this *FileLog) Close() bool {
	if this.closeState.IsClosing() {
		return true
	}
	this.closeState.AskClose()
	func() {
		defer func() {
			recover()
		}()
		this.out <- nil
	}()
	return true
}

func (this *FileLog) IsClosing() bool {
	return this.closeState.IsClosing()
}

func (this *FileLog) WaitClose() bool {
	return this.closeState.WaitClosed()
}

func (this *FileLog) Cleanup() bool {
	return this.WaitClose()
}
