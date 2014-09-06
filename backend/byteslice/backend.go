package byteslice

import (
	"oakwilson.com/p/stm"
)

type Backend struct {
	memory []byte
}

func New(memory []byte) *Backend {
	return &Backend{
		memory: memory,
	}
}

func (b Backend) ReadAt(p []byte, off int64) (n int, err error) {
	if int(off)+len(p) > len(b.memory) {
		return 0, stm.ERR_OVERRUN
	}

	copy(p, b.memory[int(off):int(off)+len(p)])

	return len(p), nil
}

func (b *Backend) WriteV(ops ...*stm.WriteOperation) error {
	for _, o := range ops {
		copy(b.memory[int(o.Offset):], o.Data)
	}

	return nil
}
