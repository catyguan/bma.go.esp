package thriftpoint

import (
	"app/cacheserver"
	"errors"
	"logger"
)

const (
	tag = "cacheserver-thriftpoint"
)

type TCacheServerImpl struct {
	service *cacheserver.CacheService
}

func NewTCacheServerImpl(s *cacheserver.CacheService) *TCacheServerImpl {
	this := new(TCacheServerImpl)
	this.service = s
	return this
}

func (this *TCacheServerImpl) CacheServerGet(treq *TCacheRequest, options map[string]string) (r *TCacheResult, err error) {
	logger.Debug(tag, "CacheServerGet(%v,%v)", treq, options)

	group := treq.GroupName
	req := cacheserver.NewGetRequest(treq.Key)
	if treq.Timeout == 0 {
		req.TimeoutMs = 5 * 1000
	} else {
		req.TimeoutMs = treq.Timeout * 1000
	}
	req.NotLoad = treq.NotLoad
	req.Trace = treq.Trace

	rep := make(chan *cacheserver.GetResult, 1)
	defer close(rep)

	err = this.service.Get(group, req, rep)
	if err != nil {
		logger.Warn(tag, "CacheServerGet fail - %s", err.Error())
		return nil, err
	}

	result := <-rep
	if result == nil {
		err = errors.New("null result return")
		logger.Warn(tag, "CacheServerGet null result return")
		return nil, err
	}

	logger.Debug(tag, "CacheServerGet %v -> %v", req, result.Done)

	tr := NewTCacheResult()
	tr.Done = result.Done
	tr.Value = result.Value
	if result.Value != nil {
		tr.Length = int32(len(result.Value))
	}
	if result.Err != nil {
		tr.ErrorA1 = result.Err.Error()
	}
	tr.Traces = result.TraceInfo

	return tr, nil
}

func (this *TCacheServerImpl) CacheServerLoad(groupName string, key string) (err error) {
	logger.Debug(tag, "CacheServerLoad(%s,%s)", groupName, key)
	err = this.service.Load(groupName, key)
	return err
}

func (this *TCacheServerImpl) CacheServerPut(groupName string, key string, value []byte, length int32) (err error) {
	logger.Debug(tag, "CacheServerPut(%s,%s)", groupName, key)
	if length > 0 && length < int32(len(value)) {
		value = value[0:length]
	}
	err = this.service.Put(groupName, key, value, 0)
	return err
}

func (this *TCacheServerImpl) CacheServerErase(groupName string, key string) (err error) {
	logger.Debug(tag, "CacheServerErase(%s,%s)", groupName, key)

	_, err = this.service.Delete(groupName, key)
	if err != nil {
		logger.Warn(tag, "CacheServerErase fail - %s", err.Error())
		return err
	}
	return nil
}
