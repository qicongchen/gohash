/*
Copyright © 2010 John Asmuth. All Rights Reserved.

Redistribution and use in source and binary forms, with or without modification,
are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list
of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this
list of conditions and the following disclaimer in the documentation and/or other
materials provided with the distribution.

3. The name of the author may not be used to endorse or promote products derived
from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY [LICENSOR] "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF
USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR
OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/

package hash

type Hasher interface {
	Hashcode() uint64
	Equals(other Hasher) bool
}

type hasherLL struct {
	item Hasher
	value interface{}
	next *hasherLL
}

func (hll *hasherLL) add(h Hasher, value interface{}) (res *hasherLL, added bool) {
	if hll == nil {
		return &hasherLL{h, value, nil}, true
	}
	if !hll.item.Equals(h) {
		hll.next, added = hll.next.add(h, value)
	} else {
		hll.value, added = value, false
	}
	res = hll
	return
}

func (hll *hasherLL) remove(h Hasher) (res *hasherLL, removed bool) {
	if hll == nil {
		return
	}
	if hll.item.Equals(h) {
		return hll.next, true
	} else {
		hll.next, removed = hll.next.remove(h)
	}
	res = hll
	return
}

func (hll *hasherLL) get(h Hasher) (res *hasherLL) {
	if hll == nil {
		return nil
	}
	if hll.item.Equals(h) {
		return hll
	}
	return hll.next.get(h)
}

type Set struct {
	bins map[uint64]*hasherLL
	count int
}

func NewSet() (hs *Set) {
	hs = new(Set)
	hs.bins = make(map[uint64]*hasherLL)
	hs.count = 0
	return
}

func (hs *Set) Keys() (out chan Hasher) {
	out = make(chan Hasher)
	go func(out chan Hasher) {
		for _, bin := range hs.bins {
			for bin != nil {
				out <- bin.item
				bin = bin.next
			}
		}
	}(out)
	return
}

func (hs *Set) Insert(h Hasher) {
	index := h.Hashcode()
	var added bool
	hs.bins[index], added = hs.bins[index].add(h, nil)
	if added {
		hs.count++
	}
}

func (hs *Set) Remove(h Hasher) {
	index := h.Hashcode()
	var removed bool
	hs.bins[index], removed = hs.bins[index].remove(h)
	if removed {
		hs.count--
	}
}

func (hs *Set) Contains(h Hasher) bool {
	index := h.Hashcode()
	return hs.bins[index].get(h) != nil
}

func (hs *Set) Size() int {
	return hs.count
}

type Map Set

type KeyValue struct {
	Key Hasher
	Value interface{}
}

func NewMap() (hs *Map) {
	hs = (*Map)(NewSet())
	return
}

func (hs *Map) Values() (out chan interface{}) {
	out = make(chan interface{})
	go func(out chan interface{}) {
		for _, bin := range hs.bins {
			for bin != nil {
				out <- bin.value
				bin = bin.next
			}
		}
	}(out)
	return
}

func (hs *Map) KeyValues() (out chan KeyValue) {
	out = make(chan KeyValue)
	go func(out chan KeyValue) {
		for _, bin := range hs.bins {
			for bin != nil {
				out <- KeyValue{bin.item, bin.value}
				bin = bin.next
			}
		}
	}(out)
	return
}

func (hs *Map) Put(h Hasher, v interface{}) {
	index := h.Hashcode()
	var added bool
	hs.bins[index], added = hs.bins[index].add(h, v)
	if added {
		hs.count++
	}
}

func (hs *Map) Get(h Hasher) (value interface{}, ok bool) {
	index := h.Hashcode()
	if hll := hs.bins[index].get(h); hll != nil {
		value = hll.value
		ok = true
	}
	return
}
