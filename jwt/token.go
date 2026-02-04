package jwt

import (
  "fmt"

  "github.com/golang-jwt/jwt/v5"
  "github.com/yoophox/emei/errs"
  "github.com/yoophox/emei/pki"
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

func (k *token) Err() error {
  return k.err
}

func (k *token) Exchange(o ...Option) JWT {
  return nil
}

func (k *token) sign() {
  if k.Token.Raw != "" {
    return
  }

  h := k.Header
  if h == nil {
    k.err = errs.Wrap(fmt.Errorf("tok have no header"), ERR_ID_JWT_NO_HEADER)
  }

  pkiid, ok := h[COMMON_HEADER_PKI_ID]
  if !ok {
    k.err = errs.Wrap(fmt.Errorf("have no pki id in header"), ERR_ID_JWT_NO_PKI_ID)
  }

  pid_ := pkiid.(uint64)
  // prikey, err := pki.GetPriKeyByID("", pid_)
  // if err != nil {
  //   return "", errs.Wrap(err, ERR_ID_JWT_NO_PKI_ID)
  // }
  // s, err := k.SignedString(prikey)

  sstr, err := k.SigningString()
  if err != nil {
    k.err = err
    return
  }
  sig, err := pki.Sign(pid_, []byte(sstr))
  if err != nil {
    k.err = err
    return
  }

  k.Valid = true
  k.Token.Raw = sstr + "." + k.EncodeSegment(sig)
}

func (k *token) Raw() string {
  return k.Token.Raw
}

func (k *token) UID() string {
  return k.GetClaim("uid")
}

func (k *token) UNmae() string {
  return k.GetClaim("uname")
}
