package hashmap

import (
	"gohash.googlecode.com/hg/hashset"
)

type Hasher hashset.Hasher
type Map hashset.Set

type KeyValue struct {
	Key Hasher
	Value interface{}
}

func (kv KeyValue) Hashcode() uint64 {
	return kv.Key.Hashcode()
}

func (kv KeyValue) Equals(other interface{}) bool {
	return kv.Key.Equals(other)
}

func New() (hs *Map) {
	hs = (*Map)(hashset.New())
	return
}

func (hs *Map) Keys() (out chan interface{}) {
	out = make(chan interface{})
	go func(out chan interface{}) {
		for kv := range hs.KeyValues() {
			out <- kv.Key
		}
	}(out)
	return
}

func (hs *Map) Values() (out chan interface{}) {
	out = make(chan interface{})
	go func(out chan interface{}) {
		for kv := range hs.KeyValues() {
			out <- kv.Value
		}
	}(out)
	return
}

func (hs *Map) KeyValues() (out chan KeyValue) {
	out = make(chan KeyValue)
	go func(out chan KeyValue) {
		for kvi := range (*hashset.Set)(hs).Keys() {
			out <- kvi.(KeyValue)
		}
	}(out)
	return
}

func (hs *Map) Put(h Hasher, v interface{}) {
	kv := KeyValue{h, v}
	(*hashset.Set)(hs).Insert(kv)
}

func (hs *Map) Get(h Hasher) (value interface{}, ok bool) {
	kvi, ok := (*hashset.Set)(hs).Get(h)
	if ok {
		value = (kvi.(KeyValue)).Value
	}
	return
}
