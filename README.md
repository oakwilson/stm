STM
===

Software Transactional Memory for Go

Overview
--------

This library implements [Software Transactional Memory](https://en.wikipedia.org/wiki/Software_transactional_memory)
in Go. Software transactional memory allows multiple producers and consumers to
interact with a single logical memory region without fear of conflict. The
implementation of this library ensures isolation between transactions, and
atomicity of commit operations.

STM has a pluggable storage mechanism, allowing the user to implement features
like journalling if the backend doesn't support atomic writes. The default
storage backend is simply a wrapper around a `[]byte` type. Anything that
allows for arbitrary reads and atomic vectorised writes can be used as a storage
mechanism. This is described with the `Backend` interface.

You can see the API documentation here:

[http://godoc.org/oakwilson.com/p/stm](http://godoc.org/oakwilson.com/p/stm)

Installation
------------

```
$ go get oakwilson.com/p/stm
```

STM Sequences
-------------

These are some ASCII graphics showing some possible sequences of STM operations.
Code is included below each diagram showing how the operation shown would be
expressed with this library.

### Non-conflicting transactions

This diagram shows the result of two interleaved transactions where neither one
conflicts with the other.

```
[start].-------.--------------------------------[end]
        `.      `.                            /|
          `-tx()-|->r(0,2)->w(4,2)->c()----->' |
                 `.                            |
                   `-tx()->r(2,2)->r(0,2)->c()-'
```

```go
t1 := m.Tx()
t2 := m.Tx()

must(t1.ReadAt(make([]byte, 2), 0))
must(t1.WriteAt([]byte{0x00, 0x01}, 4))

must(t2.ReadAt(make([]byte, 2), 2))
must(t2.ReadAt(make([]byte, 2), 0))

must(t1.Commit())
must(t2.Commit())
```

### Write blocking write

Here we have two transactions, each trying to write to the same address. This
will cause the later-committed transaction to fail.

```
[start].-------------------------[end]
       |`.                     /
       |  `-tx()->w(0,2)->c()-'
        `.
          `-tx()->w(0,2)->c()->ERROR
```

```go
t1 := m.Tx()
t2 := m.Tx()

must(t1.WriteAt([]byte{0x00, 0x01}, 0))
must(t2.WriteAt([]byte{0x00, 0x02}, 0))

must(t1.Commit())
must(t2.Commit()) // error
```

### Read blocking write

This is much the same as the above example, except that the first operation is a
read instead of a write. This still invalidates the second transaction, as it is
not allowed to modify memory that was used during another completed transaction.

```
[start].-------------------------[end]
       |`.                     /
       |  `-tx()->r(0,2)->c()-'
        `.
          `-tx()->w(0,2)->c()->ERROR
```

```go
t1 := m.Tx()
t2 := m.Tx()

must(t1.ReadAt(make([]byte, 2), 0))
must(t2.WriteAt([]byte{0x00, 0x02}, 0))

must(t1.Commit())
must(t2.Commit()) // error
```

### Write blocking read

This is the reverse of the previous example. In this case, the first transaction
writes to some memory, and the second tries to read it. There are two ways in
which the second transaction can be invalidated:

1. The first transaction completes *after* the second transaction's `read()`,
   and the second transaction's `commit()` fails.
2. The first transaction completes *before* the second transaction's `read()`,
   and the `read()` call itself fails because the transaction detects that the
   database has been modified from underneath it.

```
[start].-------------------------[end]
       |`.                     /
       |  `-tx()->w(0,2)->c()-'
        `.
          `-tx()->r(0,2)->c()->ERROR
```

This is the first failure mode:

```go
t1 := m.Tx()
t2 := m.Tx()

must(t1.WriteAt([]byte{0x00, 0x02}, 0))
must(t2.ReadAt(make([]byte, 2), 0))

must(t1.Commit())
must(t2.Commit()) // error
```

This is the second:

```go
t1 := m.Tx()
t2 := m.Tx()

must(t1.WriteAt([]byte{0x00, 0x02}, 0))
must(t1.Commit())

must(t2.ReadAt(make([]byte, 2), 0)) // error
must(t2.Commit()) // never gets here
```

Example
-------

```go
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
```

License
-------

3-clause BSD. A copy is included with the source.
