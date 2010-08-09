package hashmap

import "testing"

type StringHasher string

func (sh StringHasher) Hashcode() (hc uint64) {
	for i, c := range sh {
		hc += uint64(c) * 2 << uint64(i)
	}
	return
}

func (sh StringHasher) Equals(other interface{}) bool {
	if s, ok := other.(StringHasher); ok {
		return s == sh
	}
	return false
}

func TestMap(t *testing.T) {
	hm := New()
	hm.Put(StringHasher("john"), "A")
	hm.Put(StringHasher("stef"), "A+")
	
	if grade, ok := hm.Get(StringHasher("john")); ok {
		if grade.(string) != "A" {
			t.Fail()
		}
	} else {
		t.Fail()
	}
	
	if grade, ok := hm.Get(StringHasher("stef")); ok {
		if grade.(string) != "A+" {
			t.Fail()
		}
	} else {
		t.Fail()
	}
	
	if _, ok := hm.Get(StringHasher("no one")); ok {
		t.Fail()
	}
}
