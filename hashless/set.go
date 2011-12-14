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

//target:gohash.googlecode.com/hg/hashlessset

//Hashset backed by a left-leaning red-black tree.
package hashless

import (
	"github.com/petar/GoLLRB/llrb"
)

type Hasher interface {
	Hashcode() uint64
}
type Lesser interface {
	LessThan(other interface{}) bool
}

type HashFunc func(a interface{}) uint64
type LessFunc func(a, b interface{}) bool

func MethodLessThan(a, b interface{}) bool {
	return a.(Lesser).LessThan(b)
}

func MethodHashcode(a interface{}) uint64 {
	return a.(Hasher).Hashcode()
}

type Set struct {
	bins   map[uint64]*llrb.Tree
	lesser LessFunc
	hasher HashFunc
	count  int
}

func NewSet() (this *Set) {
	this = &Set{
		lesser: MethodLessThan,
		hasher: MethodHashcode,
		bins:   make(map[uint64]*llrb.Tree),
	}
	return
}

func NewSetFuncs(hasher HashFunc, lesser LessFunc) (this *Set) {
	this = &Set{
		lesser: lesser,
		hasher: hasher,
		bins:   make(map[uint64]*llrb.Tree),
	}
	return
}

func (this *Set) Hasher(a interface{}) uint64 {
	return this.hasher(a)
}

func (this *Set) Lesser(a, b interface{}) bool {
	return this.lesser(a, b)
}

func (this *Set) Keys() (out <-chan interface{}) {
	ch := make(chan interface{})
	out = ch
	go func(in chan<- interface{}) {
		for _, bin := range this.bins {
			for item := range bin.IterAscend() {
				in <- item
			}
		}
	}(ch)
	return
}

func (this *Set) Insert(item interface{}) {
	bin := this.bins[this.hasher(item)]
	if bin == nil {
		bin = llrb.New(llrb.LessFunc(this.lesser))
		this.bins[this.hasher(item)] = bin
	}
	if bin.ReplaceOrInsert(item) == nil {
		this.count++
	}
}

func (this *Set) Remove(key interface{}) {
	bin := this.bins[this.hasher(key)]
	if bin == nil {
		return
	}
	if bin.Delete(key) != nil {
		this.count--
	}
}

func (this *Set) Get(key interface{}) (item interface{}, ok bool) {
	bin := this.bins[this.hasher(key)]
	if bin == nil {
		return
	}
	item = bin.Get(key)
	ok = item != nil
	return
}

func (this *Set) Contains(key interface{}) bool {
	bin := this.bins[this.hasher(key)]
	if bin == nil {
		return false
	}
	return bin.Has(key)
}

func (this *Set) Size() int {
	return this.count
}
