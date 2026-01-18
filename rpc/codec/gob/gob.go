package gob

import (
  "encoding/gob"
  "io"
)

var (
  NewDec = gob.NewDecoder
  NewEnc = gob.NewEncoder
)

type gob_ struct {
  *gob.Decoder
  *gob.Encoder
}

// New ...
func New(io io.ReadWriteCloser) *gob_ {
  return &gob_{
    Decoder: gob.NewDecoder(io),
    Encoder: gob.NewEncoder(io),
  }
}
