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
		t.m.Unlock()
		return 0, ERR_EXPIRED
	}

	for _, _t := range t.m.Txs() {
		if _t == t {
			continue
		}

		for _, _o := range _t.writes {
			if conflicting(o, _o) {
				t.m.Unlock()
				return 0, ERR_OVERLAP
			}
		}
	}

	t.Lock()
	t.reads = append(t.reads, &o)
	t.Unlock()

	n, err = t.m.b.ReadAt(p, off)

	t.m.Unlock()

	if err != nil {
		return 0, err
	} else {
		return n, nil
	}
}

func (t *Tx) WriteAt(p []byte, off int64) (n int, err error) {
	o := WriteOperation{
		Offset: uint64(off),
		Length: uint64(len(p)),
		Data:   p,
	}

	t.m.Lock()

	if t.m.Version() != t.v {
		t.m.Unlock()
		return 0, ERR_EXPIRED
	}

	for _, _t := range t.m.Txs() {
		if _t == t {
			continue
		}

		for _, _o := range _t.writes {
			if conflicting(o, _o) {
				t.m.Unlock()
				return 0, ERR_OVERLAP
			}
		}

		for _, _o := range _t.reads {
			if conflicting(o, _o) {
				t.m.Unlock()
				return 0, ERR_OVERLAP
			}
		}
	}

	t.Lock()
	t.writes = append(t.writes, &o)
	t.Unlock()

	t.m.Unlock()

	return len(p), nil
}

func (t *Tx) Commit() error {
	t.m.Lock()
	defer t.m.Unlock()

	t.m.RemoveTx(t)

	if t.m.Version() != t.v {
		return ERR_EXPIRED
	}

	if err := t.m.WriteV(t.writes...); err != nil {
		return err
	}

	t.m.Increment()

	return nil
}

func (t *Tx) Abort() error {
	t.m.Lock()
	defer t.m.Unlock()

	t.m.RemoveTx(t)

	return nil
}
