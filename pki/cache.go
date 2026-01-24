package pki

import "github.com/yoophox/emei/cron"

type cryptoCache struct {
  cryptos map[string]*cryptosx
  cron    cron.Cron
}

// newCryptoCache ...
func newCryptoCache() *cryptoCache {
  return nil
}
