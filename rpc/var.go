package rpc

import (
  "io"
  "reflect"
  "sync"
)

var (
  _recvs = map[string]*rcvr{}
  _wg    *sync.WaitGroup
)

// Precompute the reflect type for error.
var (
  typeOfError        = reflect.TypeFor[error]()
  typeOfReader       = reflect.TypeFor[io.Reader]()
  typeOfWriter       = reflect.TypeFor[io.Writer]()
  typeOfReaderWriter = reflect.TypeFor[io.ReadWriter]()
)
