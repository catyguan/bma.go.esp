package sql4fileloader

import (
	"fileloader"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func safeCall() {
	time.AfterFunc(5*time.Second, func() {
		fmt.Println("os exit!!!")
		os.Exit(-1)
	})
}

func TestSQLLoader(t *testing.T) {
	safeCall()

	fn := "hello.lua"

	cfg := make(map[string]interface{})
	cfg["Type"] = "sql"
	cfg["Driver"] = "mysql"
	cfg["DataSource"] = "root:root@tcp(172.19.16.195:3306)/test_db"

	fl, err := fileloader.CommonFileLoaderFactory.Create(cfg)
	if err != nil {
		t.Error(err)
		return
	}
	bs, err2 := fl.Load(fn)
	if err2 != nil {
		t.Error(err)
		return
	}
	fmt.Println(len(bs))

	fl = nil
	for i := 0; i < 5; i++ {
		runtime.GC()
		runtime.Gosched()
	}
	time.Sleep(100 * time.Millisecond)
}
