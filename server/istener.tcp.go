package svr

import (
  "context"
  "crypto/tls"
  "net"

  "github.com/yolksys/emei/errs"
  "github.com/yolksys/emei/log"
  "github.com/yolksys/emei/pki"
)

// Listen ...
func Listen(addr string) error {
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
    }

    switch tlsc.ConnectionState().NegotiatedProtocol {
    case string(RPC_ALP_GOB):
      cc = newGobCodec(c)
      go dispatchRpc(&linkTx{cc: cc, ReadWriteCloser: c, pol: nil})
    case string(RPC_ALP_JSON):
      cc = newJsonCodec(c)
      go dispatchRpc(&linkTx{cc: cc, ReadWriteCloser: c, pol: nil})
    default:
      panic("not support alp:" + tlsc.ConnectionState().NegotiatedProtocol)
    }
  }
}
