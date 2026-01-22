package svr

import (
  "encoding/gob"
  "io"
)

var (
  NewDec = gob.NewDecoder
  NewEnc = gob.NewEncoder
)

type gobTx struct {
  *gob.Decoder
  *gob.Encoder
}

// New ...
func newGobCodec(io io.ReadWriteCloser) *gobTx {
  return &gobTx{
    Decoder: gob.NewDecoder(io),
    Encoder: gob.NewEncoder(io),
  }
}
