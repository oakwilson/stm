package main

import (
	"os"
	"sync"

	"github.com/demizer/go-elog"
	"oakwilson.com/p/stm"
	"oakwilson.com/p/stm/backend/byteslice"
)

func main() {
	log.SetLevel(log.LEVEL_DEBUG)

	m := stm.NewManager(byteslice.New(make([]byte, 12)))

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			t := m.Tx()
			b := make([]byte, 4)

			if _, err := t.ReadAt(b, 0); err != nil {
				log.Errorf("error: %s\n", err.Error())
				t.Abort()
				return
			}

			if _, err := t.WriteAt([]byte("abcd"), int64((i+1)*4)); err != nil {
				log.Errorf("error: %s\n", err.Error())
				t.Abort()
				return
			}

			if err := t.Commit(); err != nil {
				log.Errorf("error: %s\n", err.Error())
			}
		}(i)
	}

	wg.Wait()

	t := m.Tx()
	b := make([]byte, 12)

	log.Debugf("%#v\n", m)

	if _, err := t.ReadAt(b, 0); err != nil {
		log.Errorf("error: %s\n", err.Error())
		os.Exit(1)
	}

	if err := t.Commit(); err != nil {
		log.Errorf("error: %s\n", err.Error())
		os.Exit(1)
	}

	log.Debugf("finished: %x\n", b)
}
