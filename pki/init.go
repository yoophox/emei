package pki

import "github.com/yoophox/emei/cla"

func init() {
}

func init() {
  _pkiMtls = cla.Bool("pki.mtls", "mtls", false)
  //_pkiService = cla.String("pki.service", "pki service name", "")
  _pkiLocal = cla.String("pki.local", "pki local data path", "")
  _pkiLocalCertPath = cla.String("pki.local.cert", "local cert path for server", "")
  _pkiLocalKeyPath = cla.String("pki.local.key", "local private key path for server", "")

  if _pkiLocal != "" {
    _backend = newLocalBackend()
  } else {
    _backend = newRemoteBackend()
  }
}

var (
  _pkiMtls = false
  //_pkiService       = ""
  _pkiLocal         = ""
  _pkiLocalCertPath = ""
  _pkiLocalKeyPath  = ""
)
