package session

import (
  "context"
  "io"
  "net/http"

  "github.com/quic-go/quic-go"
  "github.com/quic-go/quic-go/http3"
  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/errs"
  "github.com/yolksys/emei/log"
  "github.com/yolksys/emei/pki"
  "github.com/yolksys/emei/rpc/codec"
  "github.com/yolksys/emei/rpc/errors"
)

// listenQuicSesn ...
func listenQuicSesn(addr string, ch chan<- SesnIx) error {
  tlsc, err := pki.NewServerTlsConfig()
  if err != nil {
    return err
  }
  tlsc.NextProtos = []string{string(RPC_CODEC_GOB), string(RPC_CODEC_HTTP3_JSON)}
  l, err := quic.ListenAddrEarly(addr, tlsc, &quic.Config{Allow0RTT: true})
  if err != nil {
    return errs.Wrap(err, errors.ERR_ID_RPC_QUIC_LISTEN)
  }

  go accept(l, ch)

  return nil
}

// accept ...
func accept(l *quic.EarlyListener, ch chan<- SesnIx) {
  for {
    conn, err := l.Accept(context.TODO())
    if err != nil {
      log.Event("quic accept err", err.Error())
      return
    }

    // go stream(conn, ch)
    proc := conn.ConnectionState().TLS.NegotiatedProtocol
    switch proc {
    case string(RPC_CODEC_HTTP3_JSON):
      go httpQ(conn, ch)
    default:
      go stream(conn, ch)
    }
  }
}

// stream ...
func stream(c *quic.Conn, ch chan<- SesnIx) {
  for {
    s, err := c.AcceptStream(context.TODO())
    if err != nil {
      log.Error("quic conn err", err.Error())
      return
    }

    var cc codec.CodecIx
    switch c.ConnectionState().TLS.NegotiatedProtocol {
    case string(RPC_CODEC_GOB):
      cc = codec.NewGob(s)
      go acceptSesn(cc, s, ch)

    default:
      log.Error("tls_proc", c.ConnectionState().TLS.NegotiatedProtocol)
      s.Close()
    }

  }
}

// acceptSesn ...
func acceptSesn(cc codec.CodecIx, s *quic.Stream, ch chan<- SesnIx) {
  for {
    var tja env.Tjatse
    err := cc.Decode(&tja)
    if err != nil {
      s.Close()
      log.Error("quic:decode tja", err.Error())
      return
    }

    ses := newSession(env.New(&tja), cc, s, true)
    ch <- ses
    ses.Wait()
  }
}

// ...
func httpQ(c *quic.Conn, ch chan<- SesnIx) {
  s := http3.Server{}
  s.Handler = http.HandlerFunc(serveHttp)
  s.ServeQUICConn(c)
}

// dialQuic ...
func dialQuic(addr string) (io.ReadWriteCloser, error) {
  return nil, nil
}
