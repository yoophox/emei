package jwt

import (
  "fmt"

  "github.com/golang-jwt/jwt/v5"
  "github.com/yoophox/emei/errs"
  "github.com/yoophox/emei/pki"
)

type JWT interface {
  GetClaim(key string) string
  // SetClaim(k, v string)
  Exchange(opts ...Option) JWT
  IsLegal() bool
  Err() error
  Sign() (string, error)
  Raw() string
}

// New ...
func New(o ...Option) JWT {
  opts := &options{}
  for _, f := range o {
    f(opts)
  }

  clms := &claims{}
  clms.Issuer = opts.Issuer
  clms.Subject = opts.Subject
  clms.ExpiresAt = jwt.NewNumericDate(opts.ExpiresAt)
  clms.Clms = opts.Clms

  t := &token{}
  t.Token = jwt.NewWithClaims(jwt.SigningMethodEdDSA, clms)
  // t.Header = map[string]any{}
  pkiid, err := pki.GetRandomCrpto(pki.WithED25519())
  if err != nil {
    t.err = err
    return t
  }
  t.Header[COMMON_HEADER_PKI_ID] = pkiid
  // t.Claims = clms
  t.Valid = true

  return t
}

// FromStr ...
func FromStr(jwtstr string) JWT {
  tok, err := jwt.Parse(jwtstr, func(t *jwt.Token) (any, error) {
    h := t.Header
    if h == nil {
      return nil, errs.Wrap(fmt.Errorf("tok have no header"), ERR_ID_JWT_NO_HEADER)
    }

    pkiid, ok := h[COMMON_HEADER_PKI_ID]
    if ok {
      return nil, errs.Wrap(fmt.Errorf("have no pki id in header"), ERR_ID_JWT_NO_PKI_ID)
    }

    pid_ := pkiid.(uint64)
    pubkey, err := pki.GetPubKeyByID(pid_)
    if err != nil {
      return nil, errs.Wrap(err, ERR_ID_JWT_NO_PKI_ID)
    }

    return pubkey, nil
  })

  tt_ := &token{}
  if tok == nil {
    tt_.Token = &jwt.Token{}
  } else {
    tt_.Token = tok
  }
  tt_.err = err

  return tt_
}
