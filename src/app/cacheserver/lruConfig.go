package cacheserver

import "fmt"

type LruCacheConfig struct {
	Maxsize      int32 // LRU cache maxsize
	ValidSeconds int32 // LRU cache item valid seconds, =0 mean no valid time
}

func (this *LruCacheConfig) Valid() string {
	if this.Maxsize <= 0 {
		return fmt.Sprintf("Maxsize %d invalid", this.Maxsize)
	}
	if this.ValidSeconds < 0 {
		return fmt.Sprintf("ValidSeconds %d invalid", this.ValidSeconds)
	}
	return ""
}

func (this *LruCacheConfig) NewCache() *Cache {
	cache := NewCache(this.Maxsize)
	if this.ValidSeconds > 0 {
		cache.ValidTime = int64(this.ValidSeconds)
	}
	return cache
}
