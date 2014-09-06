package stm

import (
	"errors"
	"sync"
)

var (
	ERR_OVERLAP = errors.New("overlapping operation")
	ERR_EXPIRED = errors.New("transaction has expired")
)

type Tx struct {
	sync.RWMutex

	m *Manager
	v uint64

	writes []*WriteOperation
	reads  []*ReadOperation

	laterLock sync.Mutex
	later     []*Tx

	completed bool
}

func NewTx(m *Manager) *Tx {
	return &Tx{
		m: m,
		v: m.Version(),
	}
}

func (t *Tx) ReadAt(p []byte, off int64) (n int, err error) {
	o := ReadOperation{
		Offset: uint64(off),
		Length: uint64(len(p)),
	}

	t.m.Lock()

	if t.m.Version() != t.v {
		t.laterLock.Lock()

		for _, _t := range t.later {
			for _, _o := range _t.writes {
				if overlapping(o, _o) {
					t.laterLock.Unlock()
					t.m.Unlock()

					return 0, ERR_EXPIRED
				}
			}
		}

		t.laterLock.Unlock()
	}

	t.Lock()
	t.reads = append(t.reads, &o)
	t.Unlock()

	n, err = t.m.b.ReadAt(p, off)

	t.m.Unlock()

	if err != nil {
		return 0, err
	}

	for _, _o := range t.writes {
		if !overlapping(o, _o) {
			continue
		}

		var srcLower, dstLower, l uint64

		if _o.Offset > o.Offset {
			srcLower = _o.Offset - o.Offset
			dstLower = 0
		} else {
			srcLower = 0
			dstLower = o.Offset - _o.Offset
		}

		if _o.Length > o.Length {
			l = o.Length
		} else {
			l = _o.Length
		}

		copy(p[dstLower:dstLower+l], _o.Data[srcLower:srcLower+l])
	}

	return n, nil
}

func (t *Tx) WriteAt(p []byte, off int64) (n int, err error) {
	o := WriteOperation{
		Offset: uint64(off),
		Length: uint64(len(p)),
		Data:   p,
	}

	t.m.Lock()

	t.Lock()
	t.writes = append(t.writes, &o)
	t.Unlock()

	t.m.Unlock()

	return len(p), nil
}

func (t *Tx) Commit() error {
	t.m.Lock()
	defer t.m.Unlock()

	defer t.m.RemoveTx(t)

	t.laterLock.Lock()

	for _, _t := range t.later {
		for _, o := range t.reads {
			for _, _o := range _t.writes {
				if overlapping(o, _o) {
					t.laterLock.Unlock()

					return ERR_EXPIRED
				}
			}
		}

		for _, o := range t.writes {
			for _, _o := range _t.writes {
				if overlapping(o, _o) {
					t.laterLock.Unlock()

					return ERR_OVERLAP
				}
			}

			for _, _o := range _t.reads {
				if overlapping(o, _o) {
					t.laterLock.Unlock()

					return ERR_OVERLAP
				}
			}
		}
	}

	t.laterLock.Unlock()

	if err := t.m.WriteV(t.writes...); err != nil {
		return err
	}

	t.completed = true

	t.m.Increment()

	return nil
}

func (t *Tx) Abort() error {
	t.m.Lock()
	defer t.m.Unlock()
	defer t.m.RemoveTx(t)

	return nil
}
