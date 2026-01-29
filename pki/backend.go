package pki

import (
  "encoding/pem"
)

type backend interface {
  getRandomCrypto(o ...Option) (uint64, error)
  getPriKeyByID(jwt string, id uint64) (*pem.Block, error)
  sign(id uint64, c []byte) ([]byte, error)
  getPubKeyByID(id uint64) (*pem.Block, error)
}

var _backend backend
