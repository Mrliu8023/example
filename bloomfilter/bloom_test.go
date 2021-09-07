package bloomfilter

import (
	"fmt"
	"math"
	"testing"
)

func TestFilter(t *testing.T) {
	f := NewFilter()
	fmt.Println(math.Ceil(float64(f.Cap()/8/1024/1024)))
}
