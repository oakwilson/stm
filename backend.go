package stm

import (
	"errors"
)

var (
	ERR_OVERRUN = errors.New("operation would overrun")
)

type Backend interface {
	ReadAt(p []byte, off int64) (n int, err error)
	WriteV(ops ...*WriteOperation) error
}
