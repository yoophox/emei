package env

import (
  "fmt"
  "sync"
)

type Err struct {
  Code uint16
  Err  error
}

func (e *env) Err() error {
  return e.err
}

// code 1-255 reserved
func (e *env) Errorf(code uint16, f string, args ...any) error {
  err_ := errPool.Get().(*Err)
  err_.Code = code
  err_.Err = fmt.Errorf(f, args...)
  return err_
}

func (e *Err) Error() string {
  return e.Err.Error()
}

const (
  Ok byte = iota
  InternalServerErr
  StreamServerErr
)

var errPool = &sync.Pool{
  New: func() any {
    return &Err{}
  },
}
