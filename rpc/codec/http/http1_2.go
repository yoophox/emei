package http

import (
  "fmt"
  "io"

  "golang.org/x/net/http2"
)

// newHttp1_2 ...
func newHttp1_2(io_ io.ReadWriteCloser) error {
  const preface = "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"
  b := make([]byte, len(preface))
  if _, err := io.ReadFull(io_, b); err != nil {
    return err
  }
  if string(b) != preface {
  }

  framer := http2.NewFramer(io_, io_)
  frame, err := framer.ReadFrame()
  fmt.Println(frame, err)
  return err
}
