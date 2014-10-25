package vmmsql

import (
	"database/sql"
	"golua"
	"strings"
	"time"
)

type unknowValue struct {
	Value interface{}
}

func (this *unknowValue) Scan(value interface{}) error {
	this.Value = value
	return nil
}

func FetchRow(rows *sql.Rows, ts map[string]string) (map[string]interface{}, error) {
	cols, err0 := rows.Columns()
	if err0 != nil {
		return nil, err0
	}
	desc := make([]interface{}, len(cols))
	for i, n := range cols {
		var holder interface{}
		typ, ok := ts[n]
		if ok {
			switch typ {
			case "", "string":
				holder = &sql.NullString{}
			case "int", "int32", "int64":
				holder = &sql.NullInt64{}
			case "float32", "float", "float64":
				holder = &sql.NullFloat64{}
			case "bool":
				holder = &sql.NullBool{}
			default:
				if strings.HasPrefix(typ, "time") {
					holder = &sql.NullString{}
				}
			}
		}
		if holder != nil {
			desc[i] = holder
		} else {
			desc[i] = &unknowValue{}
		}
	}
	err := rows.Scan(desc...)
	if err != nil {
		return nil, err
	}
	r := make(map[string]interface{})
	for i, n := range cols {
		typ, _ := ts[n]
		holder := desc[i]
		var v interface{}
		switch rh := holder.(type) {
		case *unknowValue:
			v = rh.Value
			if bs, ok := v.([]byte); ok {
				v = string(bs)
			}
		case *sql.NullString:
			if rh.Valid {
				switch typ {
				case "", "string":
					v = rh.String
				default:
					fm := "2006-01-02 15:04:05"
					if strings.HasPrefix(typ, "time.") {
						fm = strings.TrimPrefix(typ, "time.")
					}
					tm, err0 := time.ParseInLocation(fm, rh.String, time.Local)
					if err0 != nil {
						return nil, err0
					}
					v = golua.CreateGoTime(&tm)
				}
			}
		case *sql.NullInt64:
			if rh.Valid {
				switch typ {
				case "int", "int32":
					v = int32(rh.Int64)
				case "int64":
					v = rh.Int64
				}
			}
		case *sql.NullFloat64:
			if rh.Valid {
				switch typ {
				case "float32":
					v = float32(rh.Float64)
				case "float64":
					v = rh.Float64
				}
			}
		case *sql.NullBool:
			if rh.Valid {
				v = rh.Bool
			}
		}
		r[n] = v
	}
	return r, nil
}
