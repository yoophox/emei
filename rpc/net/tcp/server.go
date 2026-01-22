package tcp

import (
  "crypto/tls"
  "io"
  "net"

  "github.com/yolksys/emei/errs"
  "github.com/yolksys/emei/log"
  "github.com/yolksys/emei/pki"
  "github.com/yolksys/emei/rpc/errors"
)

// Listen ...
func Listen(addr string, ch chan io.ReadWriteCloser) error {
  tlsc, err := pki.NewServerTlsConfig()
  if err != nil {
    return errs.Wrap(err, errors.ERR_ID_RPC_TCP_LISTEN)
  }
  tlsc.NextProtos = []string{"h2"}
  l, err := tls.Listen("tcp", addr, tlsc)
  if err != nil {
    return errs.Wrap(err, errors.ERR_ID_RPC_TCP_LISTEN)
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
