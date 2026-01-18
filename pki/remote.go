package pki

import "encoding/pem"

type remote struct{}

// newRemoteBackend ...
func newRemoteBackend() *remote {
  return &remote{}
}

func (r *remote) getRandomCrypto(o ...Option) (uint64, error)
func (r *remote) getPriKeyByID(id uint64) (*pem.Block, error)
func (r *remote) getPubKeyByID(id uint64) (*pem.Block, error)
