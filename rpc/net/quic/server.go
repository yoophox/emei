package quic

import (
  "context"
  "io"

  "github.com/quic-go/quic-go"
  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/log"
  "github.com/yolksys/emei/pki"
  "github.com/yolksys/emei/rpc/errors"
)

// Listen ...
func Listen(addr string, ch chan io.ReadWriteCloser) error {
  l, err := quic.ListenAddrEarly(addr, pki.NewServerTlsConfig(), &quic.Config{Allow0RTT: true})
  if err != nil {
    return env.Errorf(errors.ERR_ID_RPC_QUIC_LISTEN, err)
  }

  go accept(l, ch)

  return nil
}

// accept ...
func accept(l *quic.EarlyListener, ch chan io.ReadWriteCloser) {
  for {
    conn, err := l.Accept(context.TODO())
    if err != nil {
      log.Event("quic accept err", err.Error())
      return

    }

    go stream(conn, ch)
  }
}

// stream ...
func stream(c quic.Connection, ch chan io.ReadWriteCloser) {
  for {
    s, err := c.AcceptStream(context.TODO())
    if err != nil {
      log.Event("quic conn err", err.Error())
      return
    }

    ch <- s
  }
}
