package nrpc

import (
  "errors"
  "sync"
)

var (
  _nrpc    *nrpc
  _srvPool = &sync.Pool{
    // New: newServer,
  }
  _maxSrvPerConn     = 10
  _streamEofData     = []byte{255, 255}
  _streamMaxWriteLen = 4096
  _errSvrResed       = errors.New("svrreserror")
)
