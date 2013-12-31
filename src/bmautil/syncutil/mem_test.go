package syncutil

import (
	"testing"
)

func BenchmarkAtomic(b *testing.B) {
	var m MemHolder
	for i := 0; i < b.N; i++ {
		m.Get()
	}
}
