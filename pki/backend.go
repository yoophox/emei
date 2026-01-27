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

type defaultBackend struct{}

func (d *defaultBackend) getRandomCrypto(o ...Option) (uint64, error)
func (d *defaultBackend) getPriKeyByID(jwt string, id uint64) (*pem.Block, error)
func (d *defaultBackend) getPubKeyByID(id uint64) (*pem.Block, error)
func (d *defaultBackend) sign(id uint64, c []byte) ([]byte, error)

var _backend backend
