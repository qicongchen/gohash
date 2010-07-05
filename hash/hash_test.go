package hash

import "testing"

type StringHasher string

func (sh StringHasher) Hashcode() (hc uint64) {
	for i, c := range sh {
		hc += uint64(c) * 2 << uint64(i)
	}
	return
}

func (sh StringHasher) Equals(other Hasher) bool {
	if s, ok := other.(StringHasher); ok {
		return s == sh
	}
	return false
}

func TestSet(t *testing.T) {
	hs := NewSet(20)
	hs.Insert(StringHasher("hello, world!"))
	hs.Insert(StringHasher("hello, there!"))
	hs.Insert(StringHasher("this is a sentence."))
	if !hs.Contains(StringHasher("hello, world!")) {
		t.Fail()
	}
	if !hs.Contains(StringHasher("hello, there!")) {
		t.Fail()
	}
	if !hs.Contains(StringHasher("this is a sentence.")) {
		t.Fail()
	}
	if hs.Contains(StringHasher("something else")) {
		t.Fail()
	}
	if hs.Size() != 3 {
		t.Fail()
	}
	hs.Insert(StringHasher("hello, world!"))
	if hs.Size() != 3 {
		t.Fail()
	}
	hs.Remove(StringHasher("hello, there!"))
	if hs.Contains(StringHasher("hello, there!")) {
		t.Fail()
	}
	if hs.Size() != 2 {
		t.Fail()
	}
}

func TestMap(t *testing.T) {
	hm := NewMap(20)
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
