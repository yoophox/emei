package tcp

import (
  "crypto/tls"
  "io"
  "net"

  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/log"
  "github.com/yolksys/emei/pki"
  "github.com/yolksys/emei/rpc/errors"
)

// Listem ...
func Listem(addr string, ch chan io.ReadWriteCloser) error {
  l, err := tls.Listen("tcp", addr, pki.NewServerTlsConfig())
  if err != nil {
    return env.Errorf(errors.ERR_ID_RPC_TCP_LISTEN, err)
  }

  go accept(l, ch)

  return nil
}

// accept ...
func accept(l net.Listener, ch chan io.ReadWriteCloser) {
  for {
    c, err := l.Accept()
    if err != nil {
      log.Event("tcp accept err", err.Error())
      return
    }

    ch <- c
  }
}
