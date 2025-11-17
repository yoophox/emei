package coder

type Encoder interface {
  Encode(v any) error
  Release()
}
