package codec

import (
  "encoding/gob"
  "io"
  "reflect"
)

var (
  NewDec = gob.NewDecoder
  NewEnc = gob.NewEncoder
)

type gobTx struct {
  *gob.Decoder
  *gob.Encoder
}

func (g *gobTx) DecodeTyps(typs ...reflect.Type) (vs []reflect.Value, err error) {
  return
}

func (g *gobTx) EncodeValues(vs ...reflect.Value) error {
  return nil
}

// New ...
func NewGob(io io.ReadWriteCloser) *gobTx {
  return &gobTx{
    Decoder: gob.NewDecoder(io),
    Encoder: gob.NewEncoder(io),
  }
}
