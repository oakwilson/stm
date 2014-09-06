package byteslice

import (
	"log"

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
	log.Printf("reading %d from %d", len(p), off)

	if int(off)+len(p) > len(b.memory) {
		return 0, stm.ERR_OVERRUN
	}

	copy(p, b.memory[int(off):int(off)+len(p)])

	return len(p), nil
}

func (b *Backend) WriteV(ops ...*stm.WriteOperation) error {
	for _, o := range ops {
		log.Printf("writing %x to %d", o.Data, o.Offset)

		copy(b.memory[int(o.Offset):], o.Data)
	}

	return nil
}
