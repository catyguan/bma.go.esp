package sqlloader

import (
	"app/cacheserver"
	"bmautil/sqlutil"
	"bmautil/valutil"
	"boot"
	"encoding/json"
	"errors"
	"esp/namedsql"
	"esp/shell"
	"fmt"
	"logger"
)

const (
	tag = "sqlLoader"
)

type sqlProperties struct {
	service *namedsql.NamedSQL
	// prop
	dataSource string
	query      string
	format     string // default: json, fieldName(string)
	keyType    string // default: string, int
}

func (this *sqlProperties) GetUProperties() []*shell.UProperty {
	r := make([]*shell.UProperty, 0)
	p1 := shell.NewUProperty("ds", this.dataSource, false, "namedSQL", func(v string) error {
		if this.service.CheckConfig(v) == nil {
			return errors.New(fmt.Sprintf("namedSQL[%s] not exists", v))
		}
		this.dataSource = v
		return nil
	})
	r = append(r, p1)
	p2 := shell.NewUProperty("query", this.query, false, "select query, param1=cache key", func(v string) error {
		this.query = v
		return nil
	})
	r = append(r, p2)
	p3 := shell.NewUProperty("format", this.format, true, "result format,json|fieldName", func(v string) error {
		this.format = v
		return nil
	})
	r = append(r, p3)
	p4 := shell.NewUProperty("keytype", this.keyType, true, "key type, string|int", func(v string) error {
		this.keyType = v
		return nil
	})
	r = append(r, p4)
	return r
}

func (this *sqlProperties) ToMap() map[string]interface{} {

	r := make(map[string]interface{})
	r["dataSource"] = this.dataSource
	r["query"] = this.query
	r["format"] = this.format
	r["keyType"] = this.keyType

	cfg := this.service.CheckConfig(this.dataSource)
	if cfg == nil {
		logger.Warn(tag, "invalid sql[%s]", this.dataSource)
	} else {
		m := valutil.BeanToMap(cfg)
		delete(m, "Name")
		r["inner"] = m
	}

	return r
}

func (this *sqlProperties) FromMap(vs map[string]interface{}) error {
	this.dataSource = valutil.ToString(vs["dataSource"], "")
	this.query = valutil.ToString(vs["query"], "")
	this.format = valutil.ToString(vs["format"], "")
	this.keyType = valutil.ToString(vs["keyType"], "")

	cfg := this.service.CheckConfig(this.dataSource)
	if cfg == nil {
		// create it
		if v, ok := vs["inner"]; ok {
			cfg = new(namedsql.SQLConfig)
			if valutil.ToBean(v.(map[string]interface{}), cfg) {
				logger.Info(tag, "create namedSQL[%s] for sqlLoader", this.dataSource)
				cfg.Name = this.dataSource
				err := this.service.CreateSQL(cfg)
				if err != nil {
					return err
				}
			} else {
				cfg = nil
			}
		}
	}
	if cfg == nil {
		return errors.New(fmt.Sprintf("dataSource[%s] invalid", this.dataSource))
	}
	return nil
}

type loaderProviderSQL struct {
	service *namedsql.NamedSQL
	name    string
}

func (this *loaderProviderSQL) Type() string {
	if this.name == "" {
		return "sql"
	}
	return this.name
}

func (this *loaderProviderSQL) CreateProperty() cacheserver.LoaderProperty {
	r := new(sqlProperties)
	r.service = this.service
	return r
}

func (this *loaderProviderSQL) CreateLoader(cfg *cacheserver.LoaderConfig, prop cacheserver.LoaderProperty) (cacheserver.Loader, error) {
	r := new(loaderSQL)
	r.service = this.service
	r.cache = cfg.CacheName
	r.name = cfg.LoaderName
	r.prop = prop.(*sqlProperties)
	return r, nil
}

type loaderSQL struct {
	service *namedsql.NamedSQL
	prop    *sqlProperties
	cache   string
	name    string
}

func defaultFormat(rs map[string]interface{}) (bool, []byte, error) {
	if len(rs) == 1 {
		for _, v := range rs {
			r := valutil.ToString(v, "")
			return true, []byte(r), nil
		}
		return false, nil, nil
	} else {
		return jsonFormat(rs)
	}
}

func jsonFormat(rs map[string]interface{}) (bool, []byte, error) {
	r, err := json.Marshal(rs)
	return true, r, err
}

func fieldFormat(name string, rs map[string]interface{}) (bool, []byte, error) {
	if v, ok := rs[name]; ok {
		r := valutil.ToString(v, "")
		return true, []byte(r), nil
	}
	return false, nil, errors.New("unknow field '" + name + "'")
}

func (this *loaderSQL) doLoad(service *cacheserver.CacheService, req *cacheserver.GetRequest) (err error) {
	traces := make([]string, 0)

	defer func() {
		logger.Debug(tag, "end load %d:%s/%s", req.Id, req.Name, req.Key)
		p := recover()
		if p != nil {
			err = errors.New(fmt.Sprintf("%v", p))
		}
		if err != nil {
			logger.Error(tag, "sqlLoader[%s] load error - %s", this.name, err.Error())
			traces = append(traces, err.Error())
			service.LoadEnd(this.name, this.cache, false, req.Key, nil, err, traces)
		}
	}()

	logger.Debug(tag, "do load %d:%s/%s", req.Id, req.Name, req.Key)

	ds := this.prop.dataSource
	traces = append(traces, fmt.Sprintf("sql load '%s'", ds))
	db, err1 := this.service.Get(ds)
	if err1 != nil {
		logger.Debug(tag, "get dataSource[%s] error -%s", ds, err1.Error())
		return err1
	}

	var key interface{}
	key = req.Key
	if this.prop.keyType == "int" {
		key = valutil.ToInt64(key, 0)
	}
	logger.Debug(tag, "query '%s', %v", this.prop.query, key)
	rows, err2 := db.Query(this.prop.query, key)
	if err2 != nil {
		logger.Debug(tag, "query %s fail - %s", this.prop.query, err2.Error())
		return err2
	}

	res, err3 := sqlutil.FetchMap(rows, 0, 1)
	if err3 != nil {
		logger.Debug(tag, "fetch fail - %s", err3.Error())
		return err3
	}
	traces = append(traces, fmt.Sprintf("sql load, count = %d", len(res)))

	done := false
	var val []byte
	var err4 error
	if len(res) == 0 {
		val = nil
	} else {
		switch this.prop.format {
		case "":
			done, val, err4 = defaultFormat(res[0])
		case "json":
			done, val, err4 = jsonFormat(res[0])
		default:
			done, val, err4 = fieldFormat(this.prop.format, res[0])
		}
	}

	if err4 != nil {
		logger.Debug(tag, "format '%s' fail - %s", this.prop.format, err4.Error())
		return err4
	}

	traces = append(traces, "sql load done")
	service.LoadEnd(this.name, this.cache, done, req.Key, val, nil, traces)
	return nil
}

func (this *loaderSQL) Load(service *cacheserver.CacheService, req *cacheserver.GetRequest) cacheserver.LoadTask {
	go this.doLoad(service, req)
	return nil
}

func InitSQLLoader(s *namedsql.NamedSQL, name string) {
	p := new(loaderProviderSQL)
	if s == nil {
		s = boot.ObjectFor("namedSQL").(*namedsql.NamedSQL)
	}
	if s == nil {
		panic("no namedSQL for SQLLoader!!!")
	}
	p.service = s
	cacheserver.RegLoaderProvider(p)
}
