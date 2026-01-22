package quic

import (
  "context"
  "io"

  "github.com/quic-go/quic-go"
  "github.com/yolksys/emei/errs"
  "github.com/yolksys/emei/log"
  "github.com/yolksys/emei/pki"
  "github.com/yolksys/emei/rpc/errors"
)

// Listen ...
func Listen(addr string, ch chan io.ReadWriteCloser) error {
  tlsc, err := pki.NewServerTlsConfig()
  if err != nil {
    return err
  }
  l, err := quic.ListenAddrEarly(addr, tlsc, &quic.Config{Allow0RTT: true})
  if err != nil {
    return errs.Wrap(err, errors.ERR_ID_RPC_QUIC_LISTEN)
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
func stream(c *quic.Conn, ch chan io.ReadWriteCloser) {
  for {
    s, err := c.AcceptStream(context.TODO())
    if err != nil {
      log.Event("quic conn err", err.Error())
      return
    }

    ch <- s
  }
}
