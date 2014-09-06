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

Example
-------



License
-------

3-clause BSD. A copy is included with the source.
