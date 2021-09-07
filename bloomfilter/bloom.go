package bloomfilter

import (
	"github.com/bits-and-blooms/bloom/v3"
	"sync"
)

type filter struct {
	*bloom.BloomFilter
}

var ft filter
var once sync.Once

func NewFilter() *filter {
	once.Do(func() {
		f := bloom.NewWithEstimates(5000000, 0.01)
		ft.BloomFilter = f
	})
	return &ft
}
