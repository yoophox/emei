package jwt

import "github.com/yoophox/emei/errs"

const (
  ERR_ID_JWT_NO_HEADER   errs.ErrId = "err.jwt.no.header"
  ERR_ID_JWT_NO_PKI_ID   errs.ErrId = "err.jwt.no.pki.id"
  ERR_ID_JWT_GET_PUB_KEY errs.ErrId = "err.jwt.get.pub.key"
  ERR_ID_JWT_GET_PRI_KEY errs.ErrId = "err.jwt.get.pri.key"
)
