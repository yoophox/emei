package pki

import "encoding/pem"

type cryptosx struct {
  Pri *pem.Block
  Pub *pem.Block
  CID uint64
  Typ string // rsa...
}

type cryptoClass string
