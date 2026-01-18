package tcp

import (
  "crypto/tls"
  "io"

  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/pki"
  "github.com/yolksys/emei/rpc/errors"
)

// Dial ...
func Dial(addr string) (io.ReadWriteCloser, error) {
  c, err := tls.Dial("tcp", addr, pki.NewClientTlsConfig())
  if err != nil {
    return nil, env.Errorf(errors.ERR_ID_RPC_TCP_DIAL, err)
  }

  return c, nil
}
