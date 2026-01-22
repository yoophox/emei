package jwt

import (
  "fmt"

  "github.com/golang-jwt/jwt/v5"
  "github.com/yolksys/emei/errs"
  "github.com/yolksys/emei/pki"
)

type token struct {
  *jwt.Token
  err error
}

type claims struct {
  jwt.RegisteredClaims
  Clms map[string]string `json:"clms,omitempty"`
}

func (k *token) GetClaim(key string) string {
  c := k.Claims.(*claims)
  if c == nil {
    return ""
  }
  v, _ := c.Clms[key]
  return v
}

func (k *token) SetClaim(key string, v string) {
  c := k.Claims.(*claims)
  if c == nil {
    c = &claims{}
    k.Claims = c
  }
  c.Clms[key] = v
}

func (k *token) IsLegal() bool {
  return k.Valid
}

func (k *token) ErrInfo() error {
  return k.err
}

func (k *token) Exchange(o ...Option) JWT {
  return nil
}

func (k *token) Sign() (string, error) {
  h := k.Header
  if h == nil {
    return "", errs.Wrap(fmt.Errorf("tok have no header"), ERR_ID_JWT_NO_HEADER)
  }

  pkiid, ok := h[COMMON_HEADER_PKI_ID]
  if ok {
    return "", errs.Wrap(fmt.Errorf("have no pki id in header"), ERR_ID_JWT_NO_PKI_ID)
  }

  pid_ := pkiid.(uint64)
  prikey, err := pki.GetPriKeyByID(pid_)
  if err != nil {
    return "", errs.Wrap(err, ERR_ID_JWT_NO_PKI_ID)
  }
  s, err := k.SignedString(prikey)
  if err == nil {
    k.Raw = s
  }

  return s, err
}
