package pki

import "github.com/yoophox/emei/cla"

func init() {
  _pkiMtls = cla.Bool("pki-mtls", "mtls", false)
}

var _pkiMtls = false
