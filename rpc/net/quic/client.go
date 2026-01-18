package quic

import (
  "context"
  "io"
  "sync"
  "time"

  "github.com/quic-go/quic-go"
  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/pki"
  "github.com/yolksys/emei/rpc/errors"
)

// Dial ...
func Dial(addr string) (io.ReadWriteCloser, error) {
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
      c, err = quic.DialAddrEarly(ctx, addr, pki.NewClientTlsConfig(), &quic.Config{Allow0RTT: true})
      if err != nil {
        return nil, env.Errorf(errors.ERR_ID_RPC_QUIC_DIAL_EARLY, err)
      }

      _conns[addr] = c
    }
  } else {
    _mtx.RUnlock()
  }

  s, err := c.OpenStream()
  if err == nil {
    return s, nil
  }

  c.CloseWithError(quic.ApplicationErrorCode(quic.ApplicationErrorErrorCode), "open stream error")
  _mtx.Lock()
  delete(_conns, addr)
  _mtx.Unlock()
  return nil, env.Errorf(errors.ERR_ID_RPC_QUIC_OPEN_STREAM, err)
}

var (
  _mtx   = sync.RWMutex{}
  _conns = map[string]quic.EarlyConnection{}
)
