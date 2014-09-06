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
	for i, _t := range m.txs {
		if _t == t {
			m.txs = append(m.txs[0:i], m.txs[i+1:]...)
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
