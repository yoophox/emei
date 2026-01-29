package svr

import (
  "context"
  "net/http"

  "github.com/quic-go/quic-go"
  "github.com/quic-go/quic-go/http3"
  "github.com/yoophox/emei/errs"
  "github.com/yoophox/emei/log"
  "github.com/yoophox/emei/pki"
)

// listenQuicSesn ...
func listenQuic(addr string) error {
  tlsc, err := pki.NewServerTlsConfig()
  if err != nil {
    return err
  }
  tlsc.NextProtos = []string{string(RPC_ALP_GOB), string(RPC_ALP_HTTP3), string(RPC_ALP_JSON)}
  l, err := quic.ListenAddrEarly(addr, tlsc, &quic.Config{Allow0RTT: true})
  if err != nil {
    return errs.Wrap(err, ERR_ID_RPC_QUIC_LISTEN)
  }

  go accept(l)

  return nil
}

// accept ...
func accept(l *quic.EarlyListener) {
  for {
    conn, err := l.Accept(context.TODO())
    if err != nil {
      log.Event("quic accept err", err.Error())
      return
    }

    // go stream(conn, ch)
    proc := conn.ConnectionState().TLS.NegotiatedProtocol
    switch proc {
    case string(RPC_ALP_HTTP3):
      go http3x(conn)
    case string(RPC_ALP_GOB):
      go stream(conn)
    default:
      log.Error("tls_proc", conn.ConnectionState().TLS.NegotiatedProtocol)
      conn.CloseWithError(quic.ApplicationErrorCode(quic.ApplicationErrorErrorCode), "")
    }
  }
}

// stream ...
func stream(c *quic.Conn) {
  for {
    s, err := c.AcceptStream(context.TODO())
    if err != nil {
      log.Error("accept stream", err.Error())
      return
    }

    var cc codecIx
    switch c.ConnectionState().TLS.NegotiatedProtocol {
    case string(RPC_ALP_GOB):
      cc = newGobCodec(s)
      go dispatchRpc(&linkTx{cc: cc, Conn: wrapQuicConn(s, c), pol: nil})
    case string(RPC_ALP_JSON):
      cc = newJsonCodec(s)
      go dispatchRpc(&linkTx{cc: cc, Conn: wrapQuicConn(s, c), pol: nil})
    default:
      panic("not support alp:" + c.ConnectionState().TLS.NegotiatedProtocol)
    }

  }
}

func http3x(c *quic.Conn) {
  s := http3.Server{}
  s.Handler = http.HandlerFunc(serveHttp)
  s.ServeQUICConn(c)
}
