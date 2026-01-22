package tcp

import (
  "crypto/tls"
  "io"

  "github.com/yolksys/emei/errs"
  "github.com/yolksys/emei/pki"
  "github.com/yolksys/emei/rpc/errors"
)

// Dial ...
func Dial(addr string) (io.ReadWriteCloser, error) {
  tlsc, err := pki.NewClientTlsConfig()
  if err != nil {
    return nil, err
  }
  c, err := tls.Dial("tcp", addr, tlsc)
  if err != nil {
    return nil, errs.Wrap(err, errors.ERR_ID_RPC_TCP_DIAL)
  }

  return c, nil
}
