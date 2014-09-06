package main

import (
	"log"

	"oakwilson.com/p/stm"
	"oakwilson.com/p/stm/backend/byteslice"
)

func must(args ...interface{}) {
	if args[len(args)-1] != nil {
		panic(args[len(args)-1])
	}
}

func main() {
	b := make([]byte, 10)

	log.Printf("before: %x", b)

	m := stm.NewManager(byteslice.New(b))

	t1 := m.Tx()
	t2 := m.Tx()

	must(t1.WriteAt([]byte("xx"), 0))
	must(t2.WriteAt([]byte("yy"), 2))

	must(t1.Commit())
	must(t2.Commit())

	log.Printf("after: %x", b)
}
