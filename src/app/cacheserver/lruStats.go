package cacheserver

import (
	"bmautil/valutil"
	"bytes"
	"fmt"
)

type CacheStats struct {
	Gets      int64
	CacheHits int64
}

type LruCacheStats struct {
	CacheStats
	Size       int32
	MaxSize    int32
	MaxCollide int
	MaxRefill  int
	TotalUse   uint64
}

func (this *LruCacheStats) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	this.BuildString(buf)
	return buf.String()
}

func (this *LruCacheStats) BuildString(buf *bytes.Buffer) {
	var per float64

	buf.WriteString(fmt.Sprintf("MaxSize=%d,", this.MaxSize))

	per = 0
	if this.MaxSize > 0 {
		per = float64(this.Size*100) / float64(this.MaxSize)
	}
	buf.WriteString(fmt.Sprintf("Size=%d(%.2f", this.Size, per))
	buf.WriteString("%),")

	buf.WriteString(fmt.Sprintf("Gets=%d,", this.Gets))

	per = 0
	if this.Gets > 0 {
		per = float64(this.CacheHits*100) / float64(this.Gets)
	}
	buf.WriteString(fmt.Sprintf("CacheHits=%d(%.2f", this.CacheHits, per))
	buf.WriteString("%),")

	tsize := valutil.MakeSizeString(this.TotalUse)

	buf.WriteString(fmt.Sprintf("MaxCollide=%d,", this.MaxCollide))

	buf.WriteString(fmt.Sprintf("MaxRefill=%d,", this.MaxRefill))

	buf.WriteString(fmt.Sprintf("TotalUse=%s", tsize))

}

func (this *LruCacheStats) CopyLruCacheState(s *CacheStats, cache *Cache) {
	this.CacheStats = *s
	this.MaxCollide = cache.MaxCollide
	this.MaxRefill = cache.MaxRefill
	this.MaxSize = cache.MaxSize()
	this.TotalUse = cache.TotalUse()
}
