package coder

import "encoding/json"

type Decoder interface {
  Decode(v any) error
  Release()
}

type DefaultDecoder struct{}

func (d *DefaultDecoder) Decode(v any) error {
  return nil
}

func (d *DefaultDecoder) Release() {
}

type JsonStrArrDecoder struct {
  a []string
  i int
}

func (d *JsonStrArrDecoder) Decode(v any) error {
  defer func() {
    d.i++
  }()
  return json.Unmarshal([]byte(d.a[d.i]), v)
}

func (d *JsonStrArrDecoder) Release() {
}
