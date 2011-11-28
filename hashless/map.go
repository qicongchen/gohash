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

//target:gohash.googlecode.com/hg/hashlessmap

//Hashmap backed by a left-leaning red-black tree.
package hashless

type KeyValue struct {
	Key, Value interface{}
}

type Map Set

func NewMap() *Map {
	return NewMapFuncs(MethodHashcode, MethodLessThan)
}

func NewMapFuncs(hasher HashFunc, lesser LessFunc) (me *Map) {

	kvhasher := func(kv interface{}) uint64 {
		return hasher(kv.(KeyValue).Key)
	}
	kvlesser := func(kv1, kv2 interface{}) bool {
		return lesser(kv1.(KeyValue).Key, kv2.(KeyValue).Key)
	}

	me = (*Map)(NewSetFuncs(kvhasher, kvlesser))
	return
}

func (me *Map) Keys() (out <-chan interface{}) {
	ch := make(chan interface{})
	out = ch
	go func(in chan<- interface{}) {
		for kv := range me.KeyValues() {
			in <- kv.Key
		}
		close(in)
	}(ch)
	return
}

func (me *Map) Values() (out <-chan interface{}) {
	ch := make(chan interface{})
	out = ch
	go func(in chan<- interface{}) {
		for kv := range me.KeyValues() {
			in <- kv.Value
		}
		close(in)
	}(ch)
	return
}

func (me *Map) KeyValues() (out <-chan KeyValue) {
	ch := make(chan KeyValue)
	out = ch
	go func(in chan<- KeyValue) {
		for kvi := range (*Set)(me).Keys() {
			in <- kvi.(KeyValue)
		}
		close(in)
	}(ch)
	return
}

func (me *Map) Put(k interface{}, v interface{}) {
	kv := KeyValue{k, v}
	(*Set)(me).Insert(kv)
}

func (me *Map) Get(k interface{}) (value interface{}, ok bool) {
	kvi, ok := (*Set)(me).Get(KeyValue{k, nil})
	if ok {
		value = (kvi.(KeyValue)).Value
	}
	return
}

func (me *Map) Size() int {
	return (*Set)(me).Size()
}
