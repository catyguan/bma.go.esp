package eloader

import (
	"app/cacheserver"
	"logger"
	"os"
	"testing"
	"time"
)

func TestLoaderCache(t *testing.T) {
	c := new(LoaderCache)
	c.InitCache(nil, "test")
	c.config = new(cacheConfig)
	c.config.Maxsize = 1024

	lcfg := new(LoaderConfig)
	lcfg.Type = "none"

	if err := c.Start(); err != nil {
		t.Error("start", err)
	} else {
		time.AfterFunc(2*time.Second, func() {
			os.Exit(-99)
		})

		err2 := c.AddLoader(lcfg)
		if err2 != nil {
			t.Error("addLoader", err2)
		}

		req := cacheserver.NewGetRequest("key1")
		req.TimeoutMs = 1000
		req.Trace = true
		rep := make(chan *cacheserver.GetResult, 1)
		c.Get(req, rep)
		logger.Info("TEST", "WAIT")
		r := <-rep
		logger.Info("TEST", "key1 = %v", r)
	}
	c.Stop()
	time.Sleep(1 * time.Millisecond)
}
