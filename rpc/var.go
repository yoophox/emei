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
  typeOfReader       = reflect.TypeFor[io.ReadCloser]()
  typeOfWriter       = reflect.TypeFor[io.WriteCloser]()
  typeOfReaderWriter = reflect.TypeFor[io.ReadWriteCloser]()
  typeOfWebsock      = reflect.TypeFor[WebSock]()
  typeOfUpFiles      = reflect.TypeFor[UpFiles]()
)
