package hashless

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

func (sh StringHasher) LessThan(other interface{}) (less bool) {
	osh := other.(StringHasher)
	if len(sh) < len(osh) {
		less = true
		return
	}
	for i, _ := range sh {
		if sh[i] < osh[i] {
			less = true
			return
		}
	}
	less = false
	return
}
