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
	bins []*hasherLL
	count int
}

func NewSet(capacity int) (hs *Set) {
	hs = new(Set)
	hs.bins = make([]*hasherLL, capacity)
	hs.count = 0
	return
}

func (hs *Set) Insert(h Hasher) {
	index := h.Hashcode() % uint64(len(hs.bins))
	var added bool
	hs.bins[index], added = hs.bins[index].add(h, nil)
	if added {
		hs.count++
	}
}

func (hs *Set) Remove(h Hasher) {
	index := h.Hashcode() % uint64(len(hs.bins))
	var removed bool
	hs.bins[index], removed = hs.bins[index].remove(h)
	if removed {
		hs.count--
	}
}

func (hs *Set) Contains(h Hasher) bool {
	index := h.Hashcode() % uint64(len(hs.bins))
	return hs.bins[index].get(h) != nil
}

func (hs *Set) Size() int {
	return hs.count
}

type Map Set
type MapP *Map

func NewMap(capacity int) (hs *Map) {
	hs = (*Map)(NewSet(capacity))
	return
}

func (hs *Map) Put(h Hasher, v interface{}) {
	index := h.Hashcode() % uint64(len(hs.bins))
	var added bool
	hs.bins[index], added = hs.bins[index].add(h, v)
	if added {
		hs.count++
	}
}

func (hs *Map) Get(h Hasher) (value interface{}, ok bool) {
	index := h.Hashcode() % uint64(len(hs.bins))
	if hll := hs.bins[index].get(h); hll != nil {
		value = hll.value
		ok = true
	}
	return
}
