package pki

type cipherTx struct {
  Pri any
  Pub any
  CID uint64
  Typ string // rsa...
}

type cipherType string
