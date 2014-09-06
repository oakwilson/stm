package stm

type WriteOperation struct {
	Offset uint64
	Length uint64
	Data   []byte
}

type ReadOperation struct {
	Offset uint64
	Length uint64
}

type boundable interface {
	bounds() (uint64, uint64)
}

func overlapping(a, b boundable) bool {
	a1, a2 := a.bounds()
	b1, b2 := b.bounds()

	return !(a1 >= b2 || b1 >= a2)
}

func (w WriteOperation) bounds() (uint64, uint64) {
	return w.Offset, w.Offset + w.Length
}

func (r ReadOperation) bounds() (uint64, uint64) {
	return r.Offset, r.Offset + r.Length
}
