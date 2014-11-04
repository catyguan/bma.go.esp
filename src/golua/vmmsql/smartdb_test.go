package vmmsql

import (
	"fmt"
	"testing"
	"time"
)

func T2estSmartDB(t *testing.T) {
	if true {
		safeCall()
		doSmartDB()
	}
}

func doSmartDB() {
	if true {
		dbi := new(dbInfo)
		dbi.Name = "test"
		dbi.Driver = "mysql"
		dbi.DataSource = "root:root@tcp(172.19.16.195:3306)/test_db"

		o, _ := createSmartDB()
		sdb := o.(*smartDB)

		sdb.Add(dbi, true)

		time.Sleep(1000 * time.Millisecond)

		sdbi := sdb.Select("test", false)
		fmt.Println("select => ", sdbi)
	}
}
