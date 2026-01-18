package pki

import "github.com/yolksys/emei/cron"

type cryptoCache struct {
  cryptos map[string]*cryptosx
  cron    cron.Cron
}

// newCryptoCache ...
func newCryptoCache() *cryptoCache {
  return nil
}
