package net

import (
  "io"

  "github.com/yolksys/emei/rpc/net/quic"
  "github.com/yolksys/emei/rpc/net/tcp"
)

var (
  // return io.writereadcloser
  DialQuic = ""
  DialTcp  = ""
)

// Listen ...
// @nets: tcp/all
func Listen(addr string, nets ...string) (chan io.ReadWriteCloser, error) {
  ch := make(chan io.ReadWriteCloser, 1000)

  if nets == nil || nets[0] == "all" {
    err := quic.Listen(addr, ch)
    if err != nil {
      return nil, err
    }
  }

  if nets != nil && (nets[0] == "all" || nets[0] == "tcp") {
    err := quic.Listen(addr, ch)
    if err != nil {
      return nil, err
    }
  }

  return ch, nil
}

// Dial ...
func Dial(addr string) (io.ReadWriteCloser, error) {
  c, err := quic.Dial(addr)
  if err == nil {
    return c, nil
  }

  return tcp.Dial(addr)
}
