package svr

import (
  "encoding/json"
  "io"
)

type jsonTx struct {
  json.Encoder
  json.Decoder
}

// newJsonCodec ...
func newJsonCodec(io_ io.ReadWriter) codecIx {
  return &jsonTx{*json.NewEncoder(io_), *json.NewDecoder(io_)}
}
