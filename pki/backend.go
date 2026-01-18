package pki

import (
  "encoding/pem"

  "github.com/yolksys/emei/cla"
)

type backend interface {
  getRandomCrypto(o ...Option) (uint64, error)
  getPriKeyByID(id uint64) (*pem.Block, error)
  getPubKeyByID(id uint64) (*pem.Block, error)
}

func init() {
  _pkiService = cla.String("pki", "pki service name", "")
  _pkiLocal = cla.String("pki-local", "pki local data path", "")

  if _pkiService != "" {
    _backend = newRemoteBackend()
  } else if _pkiLocal != "" {
    _backend = newLocalBackend()
  } else {
    _backend = &defaultBackend{}
  }
}

type defaultBackend struct{}

func (d *defaultBackend) getRandomCrypto(o ...Option) (uint64, error)
func (d *defaultBackend) getPriKeyByID(id uint64) (*pem.Block, error)
func (d *defaultBackend) getPubKeyByID(id uint64) (*pem.Block, error)

var (
  _pkiService = ""
  _pkiLocal   = ""

  _backend backend
)
