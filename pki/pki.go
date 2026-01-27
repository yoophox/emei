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
  var c tls.Certificate
  var err error

  if _pkiLocalCertPath != "" && _pkiLocalKeyPath != "" {
    c, err = tls.LoadX509KeyPair(_pkiLocalCertPath, _pkiLocalKeyPath)
  } else {
    cert, key, _, err := KeyPairWithPin(4096)
    if err != nil {
      return nil, err
    }
    c, err = tls.X509KeyPair(cert, key)
  }

  if err != nil {
    return nil, err
  }

  return &tls.Config{
    Certificates: []tls.Certificate{c},
  }, nil
}

// GetPriKeyByID ...
func GetPriKeyByID(jwt string, id uint64) (any, error) {
  return _backend.getPriKeyByID(jwt, id)
}

// Sign ...
func Sign(id uint64, c []byte) ([]byte, error) {
  return _backend.sign(id, c)
}

// GetPubKeyByID ...
func GetPubKeyByID(id uint64) (any, error) {
  return _backend.getPubKeyByID(id)
}

// GetRandomCrpto ...
// @return: crypto ID
func GetRandomCrpto(opts ...Option) (uint64, error) {
  return _backend.getRandomCrypto(opts...)
}

// NewCertAndKey
func NewCertAndKey(o ...Option) {
}
