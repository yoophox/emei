package pki

import "crypto/tls"

// NewClientTlsConfig ...
func NewClientTlsConfig() *tls.Config {
  return &tls.Config{
    InsecureSkipVerify: true,
  }
}
