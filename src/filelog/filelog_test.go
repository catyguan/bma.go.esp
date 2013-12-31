package filelog

import (
	"testing"
)

func TestFileLog(t *testing.T) {
	fn := "../../test/test.log"
	flog := NewFileLog(fn, 16)
	err := flog.Open()
	if err != nil {
		t.Error(err)
		return
	}

	flog.Println("hello world")
	flog.Printf("%s = %d\n", "my ago", 25)

	flog.Close()
	flog.WaitClose()
}
