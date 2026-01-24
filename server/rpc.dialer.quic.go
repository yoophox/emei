package svr

import (
  "context"
  "net"
  "sync"
  "time"

  "github.com/quic-go/quic-go"
  "github.com/yoophox/emei/errs"
  "github.com/yoophox/emei/pki"
)

// Dial ...
func dialQuic(addr string) (net.Conn, error) {
  _mtx.RLock()
  c, ok := _conns[addr]
  if !ok {
    _mtx.RUnlock()
    _mtx.Lock()
    defer _mtx.Unlock()
    c, ok = _conns[addr]
    if !ok {
      ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // 3s handshake timeout
      defer cancel()
      var err error
      tlsc, err := pki.NewClientTlsConfig()
      if err != nil {
        return nil, err
      }
      tlsc.NextProtos = []string{string(RPC_ALP_GOB)}
      c, err = quic.DialAddrEarly(ctx, addr, tlsc, &quic.Config{Allow0RTT: true})
      if err != nil {
        return nil, errs.Wrap(err, ERR_ID_RPC_QUIC_DIAL_EARLY)
      }

      _conns[addr] = c
    }
  } else {
    _mtx.RUnlock()
  }

  s, err := c.OpenStream()
  if err == nil {
    return wrapQuicConn(s, c), nil
  }

  c.CloseWithError(quic.ApplicationErrorCode(quic.ApplicationErrorErrorCode), "open stream error")
  _mtx.Lock()
  delete(_conns, addr)
  _mtx.Unlock()
  return nil, errs.Wrap(err, ERR_ID_RPC_QUIC_OPEN_STREAM)
}

var (
  _mtx   = sync.RWMutex{}
  _conns = map[string]*quic.Conn{}
)
