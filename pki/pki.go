package pki

import (
  "crypto/tls"
  "crypto/x509"
)

func VerifyPeerCertificate(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
  return nil
}

// NewClientTlsConfig ...
func NewClientTlsConfig() (*tls.Config, error) {
  return &tls.Config{
    VerifyPeerCertificate: VerifyPeerCertificate,
    // InsecureSkipVerify: true,
  }, nil
}

// NewServerTlsConfig ...
func NewServerTlsConfig() (*tls.Config, error) {
  cert, key, _, err := KeyPairWithPin(4096)
  if err != nil {
    return nil, err
  }
  c, err := tls.X509KeyPair(cert, key)
  if err != nil {
    return nil, err
  }

  return &tls.Config{
    Certificates: []tls.Certificate{c},
  }, nil
}

// GetPriKeyByID ...
func GetPriKeyByID(id uint64) (any, error) {
  return nil, nil
}

// GetPubKeyByID ...
func GetPubKeyByID(id uint64) (any, error) {
  return nil, nil
}

// GetRandomCrpto ...
// @return: crypto ID
func GetRandomCrpto(opts ...Option) (uint64, error) {
  return 0, nil
}

// NewCertAndKey
func NewCertAndKey(o ...Option) {
}
