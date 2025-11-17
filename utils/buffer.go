package utils

import (
  "bytes"
  "sync"
)

type buffer struct {
  *bytes.Buffer
  isR bool
}

func (b *buffer) Close() error {
  if b.isR {
    return nil
  }
  bufferPool.Put(b)
  return nil
}

// NewRbuffer ...
func NewRbuffer(b []byte) *buffer {
  return &buffer{bytes.NewBuffer(b), true}
}

// NewWBuffer ...
func NewWBuffer() *buffer {
  b := bufferPool.Get().(*buffer)
  b.Reset()
  return b
}

var bufferPool = sync.Pool{
  New: func() any {
    return &buffer{
      bytes.NewBuffer(make([]byte, 0)),
      false,
    }
  },
}
