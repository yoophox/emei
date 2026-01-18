package jwt

import "time"

type options struct {
  // the `iss` (Issuer) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.1
  Issuer string

  // the `sub` (Subject) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.2
  Subject string

  // the `aud` (Audience) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.3
  Audience []string

  // the `exp` (Expiration Time) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.4
  ExpiresAt time.Time

  // the `nbf` (Not Before) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.5
  NotBefore time.Time

  // the `iat` (Issued At) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.6
  IssuedAt time.Time

  // the `jti` (JWT ID) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.7
  ID string

  Clms map[string]string
}

type Option func(*options)

// WithClaims ...
func WithClaims(kv ...string) Option {
  return func(o *options) {
    l := len(kv) & 0b1111111111111110
    if o.Clms == nil {
      o.Clms = map[string]string{}
    }
    for i := range l {
      o.Clms[kv[i]] = kv[i+1]
      i += 2
    }
  }
}

// WithID ...
func WithID(id string) Option {
  return func(o *options) {
    o.ID = id
  }
}

// WithSubject ...
func WithSubject(sub string) Option {
  return func(o *options) {
    o.Subject = sub
  }
}

// WithIssuer ...
func WithIssuer(iss string) Option {
  return func(o *options) {
    o.Issuer = iss
  }
}

// WithExpireTime ...
func WithExpireTime(etm time.Time) Option {
  return func(o *options) {
    o.ExpiresAt = etm
  }
}
