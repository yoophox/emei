package pki

import "github.com/yoophox/emei/flag"

func init() {
}

func init() {
  fs_ := flag.NewFlagSet("pki")
  pmtlsP := fs_.Bool("pki.mtls", false, "if use mtls")
  pLocalP := fs_.String("pki.local", "", "pki local data path")
  pLCertPathP := fs_.String("pki.local.cert", "", "local cert path")
  pLKeyPathP := fs_.String("pki.local.key", "", "local key path")
  err := fs_.Parse()
  if err == flag.ErrHelp {
    return
  }

  //_pkiService = cla.String("pki.service", "pki service name", "")
  _pkiMtls = *pmtlsP
  _pkiLocal = *pLocalP
  _pkiLocalCertPath = *pLCertPathP
  _pkiLocalKeyPath = *pLKeyPathP

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
