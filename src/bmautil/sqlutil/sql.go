package sqlutil

import (
	"database/sql"
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

func FetchRow(rows *sql.Rows) (map[string]interface{}, error) {
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
		r[fns[i]] = data[i]
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
