package sqlutil

import (
	"bmautil/valutil"
	"database/sql"
	"time"
)

type SQLAction func(db *sql.DB) error
type RowScan func(rows *sql.Rows) error

func ExecuteAction(callback func(int64), sqlstr string, args ...interface{}) SQLAction {
	return func(db *sql.DB) error {
		res, err := db.Exec(sqlstr, args...)
		if err != nil {
			return err
		}
		if callback != nil {
			r, err := res.RowsAffected()
			if err != nil {
				return err
			}
			callback(r)
		}
		return nil
	}
}

func InsertIdAction(callback func(int64), sqlstr string, args ...interface{}) SQLAction {
	return func(db *sql.DB) error {
		res, err := db.Exec(sqlstr, args...)
		if err != nil {
			return err
		}
		if callback != nil {
			r, err := res.LastInsertId()
			if err != nil {
				return err
			}
			callback(r)
		}
		return nil
	}
}

func QueryAction(callback RowScan, sqlstr string, args ...interface{}) SQLAction {
	return func(db *sql.DB) error {
		rows, err := db.Query(sqlstr, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		if callback != nil {
			callback(rows)
		}
		return nil
	}
}

func FetchRow(rows *sql.Rows, desc map[string]string) (map[string]interface{}, error) {
	if !rows.Next() {
		return nil, nil
	}

	fns, err0 := rows.Columns()
	if err0 != nil {
		return nil, err0
	}
	sz := len(fns)
	r := make(map[string]interface{})
	data := make([]interface{}, sz)
	scaninterfaces := make([]interface{}, sz)
	for i := 0; i < sz; i++ {
		scaninterfaces[i] = &data[i]
	}

	err2 := rows.Scan(scaninterfaces...)
	if err2 != nil {
		return nil, err2
	}
	for i := 0; i < sz; i++ {
		v := data[i]
		k := fns[i]
		if v == nil {
			r[k] = nil
			continue
		}
		vt := ""
		if desc != nil {
			vt = desc[k]
		}
		if bs, ok := v.([]byte); ok {
			if vt != "bytes" {
				v = string(bs)
			} else {
				r[k] = bs
				continue
			}
		}
		switch vt {
		case "", "string":
			v = valutil.ToString(v, "")
		case "int":
			v = valutil.ToInt(v, 0)
		case "int32":
			v = valutil.ToInt32(v, 0)
		case "int64":
			v = valutil.ToInt64(v, 0)
		case "float32":
			v = valutil.ToFloat32(v, 0)
		case "float", "float64":
			v = valutil.ToFloat64(v, 0)
		case "bool":
			v = valutil.ToBool(v, false)
		case "time":
			fm := "2006-01-02 15:04:05"
			if desc != nil {
				if nfm, ok := desc[k+".format"]; ok {
					fm = nfm
				}
			}
			if tm, ok := v.(time.Time); ok {
				v = tm
			} else {
				tmp := valutil.ToString(v, "")
				v, err0 = time.ParseInLocation(fm, tmp, time.Local)
				if err0 != nil {
					v = time.Unix(0, 0)
				}
			}
		}
		r[k] = v
	}
	return r, nil
}

func FetchMap(rows *sql.Rows, spos, count int) ([]map[string]interface{}, error) {
	fns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	sz := len(fns)
	r := make([]map[string]interface{}, 0)
	data := make([]interface{}, sz)
	scaninterfaces := make([]interface{}, sz)
	for i := 0; i < sz; i++ {
		scaninterfaces[i] = &data[i]
	}
	pos := 0
	for rows.Next() {
		pos++
		if pos <= spos {
			continue
		}
		rows.Scan(scaninterfaces...)
		m := make(map[string]interface{})
		for i := 0; i < sz; i++ {
			m[fns[i]] = data[i]
		}
		r = append(r, m)
		if len(r) >= count {
			break
		}
	}
	return r, nil
}
