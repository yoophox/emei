package svr

import (
  "context"
  "crypto/tls"
  "net"
  "net/http"

  "github.com/yoophox/emei/errs"
  "github.com/yoophox/emei/log"
  "github.com/yoophox/emei/pki"
  "golang.org/x/net/http2"
)

// Listen ...
func listenTcp(addr string) error {
  tlsc, err := pki.NewServerTlsConfig()
  if err != nil {
    return errs.Wrap(err, ERR_ID_RPC_TCP_LISTEN)
  }
  tlsc.NextProtos = []string{string(RPC_ALP_HTTP2), string(RPC_ALP_GOB), string(RPC_ALP_JSON)}
  l, err := tls.Listen("tcp", addr, tlsc)
  if err != nil {
    return errs.Wrap(err, ERR_ID_RPC_TCP_LISTEN)
  }

  go acceptTcp(l)

  return nil
}

// accept ...
func acceptTcp(l net.Listener) {
  for {
    c, err := l.Accept()
    if err != nil {
      log.Event("tcp accept err", err.Error())
      return
    }

    tlsc := c.(*tls.Conn)

    var cc codecIx

    if err := tlsc.HandshakeContext(context.Background()); err != nil {
      // panic("tcp handshakecontext err: " + err.Error())
    }

    switch tlsc.ConnectionState().NegotiatedProtocol {
    case string(RPC_ALP_GOB):
      cc = newGobCodec(c)
      go dispatchRpc(&linkTx{cc: cc, Conn: c, pol: nil})
    case string(RPC_ALP_JSON):
      cc = newJsonCodec(c)
      go dispatchRpc(&linkTx{cc: cc, Conn: c, pol: nil})
    case string(RPC_ALP_HTTP2):
      serveh(c)
    default:
      log.Event("not support alp", tlsc.ConnectionState().NegotiatedProtocol)
    }
  }
}

// serveHttp ...
func serveh(c net.Conn) {
  opt := http2.ServeConnOpts{
    Context: context.Background(),
    Handler: http.HandlerFunc(serveHttp),
  }
  s := http2.Server{}
  s.ServeConn(c, &opt)
}
