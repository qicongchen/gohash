//target:gohash.googlecode.com/hg/hashutil
package hashutil

import "unsafe"

func HashFloat64(f float64) (hash uint64) {
	hash = *(*uint64)(unsafe.Pointer(&f))
	return
}
func HashFloat64Slice(fs []float64) (hash uint64) {
	for i, f := range fs {
		hf := HashFloat64(f)
		hash += hf << uint(i)
		hash += hf >> uint(64-i)
	}
	return
}
