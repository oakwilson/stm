package stm

import (
	"sync"
)

type Manager struct {
	sync.RWMutex

	txs []*Tx
	v   uint64

	b Backend
}

func NewManager(b Backend) *Manager {
	return &Manager{
		b: b,
	}
}

func (m *Manager) Version() uint64 {
	return m.v
}

func (m *Manager) Increment() {
	m.v++
}

func (m *Manager) Tx() *Tx {
	m.Lock()
	defer m.Unlock()

	t := NewTx(m)

	m.AddTx(t)

	return t
}

func (m *Manager) AddTx(t *Tx) {
	m.txs = append(m.txs, t)
}

func (m *Manager) RemoveTx(t *Tx) {
	i := 0
	for {
		_t := m.txs[i]

		if _t == t {
			m.txs[i] = m.txs[len(m.txs)-1]
			m.txs[len(m.txs)-1] = nil
			m.txs = m.txs[0 : len(m.txs)-1]
		} else {
			i++

			if t.completed {
				_t.laterLock.Lock()
				_t.later = append(_t.later, t)
				_t.laterLock.Unlock()
			}
		}

		if i >= len(m.txs) {
			break
		}
	}
}

func (m *Manager) Txs() []*Tx {
	return m.txs
}

func (m Manager) WriteV(ops ...*WriteOperation) error {
	return m.b.WriteV(ops...)
}
