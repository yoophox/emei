package svr

import (
  "net/http"
  "time"

  "github.com/yoophox/emei/cfg"
)

// NewCookie ...
func NewCookie(opts ...optionCok) *http.Cookie {
  o := &optionCokTx{
    path:     "/" + cfg.Service + "/cookies",
    secure:   true,
    httpOnly: true,
    maxAge:   int(time.Hour) * 25,
    sameSite: http.SameSiteStrictMode,
  }
  for _, v := range opts {
    v(o)
  }

  cok := &http.Cookie{
    Name:     o.name,
    Value:    o.value,
    Path:     o.path,
    Secure:   o.secure,
    HttpOnly: o.httpOnly,
    SameSite: o.sameSite,

    Quoted: o.quoted,
  }
  return cok
}

// WithoutHttpOnlyCookie ...
func WithoutHttpOnlyCookie() optionCok {
  return func(o *optionCokTx) {
    o.httpOnly = false
  }
}

// WithoutSecureCookie ...
func WithoutSecureCookie() optionCok {
  return func(o *optionCokTx) {
    o.secure = false
  }
}

// WithQuotedCookie ...
func WithQuotedCookie() optionCok {
  return func(o *optionCokTx) {
    o.quoted = true
  }
}

// WithJwtCookie ...
func WithJwtCookie(j string) optionCok {
  return func(o *optionCokTx) {
    o.name = RPC_JWT_COOKIES_NAME
    o.value = j
  }
}

func WithDelJwtCookie() optionCok {
  return func(o *optionCokTx) {
    o.name = RPC_JWT_COOKIES_NAME
    o.maxAge = -1
  }
}

func WithTIDCookie(v string) optionCok {
  return func(o *optionCokTx) {
    o.name = RPC_TID_COOKIES_NAME
    o.value = v
  }
}

func WithDelTIDCookie() optionCok {
  return func(o *optionCokTx) {
    o.name = RPC_TID_COOKIES_NAME
    o.maxAge = -1
  }
}

// WithKVCookie ...
func WithKVCookie(k, v string) optionCok {
  return func(o *optionCokTx) {
    o.name = k
    o.value = v
  }
}

func WithDelCookie(k string) optionCok {
  return func(o *optionCokTx) {
    o.name = k
    o.maxAge = -1
  }
}

// WithMaxAgeCookie ...
func WithMaxAgeCookie(s int) optionCok {
  return func(o *optionCokTx) {
    o.maxAge = s
  }
}

// WithAllSiteCookie ...
func WithAllSiteCookie() optionCok {
  return func(o *optionCokTx) {
    o.sameSite = http.SameSiteNoneMode
  }
}

// AddJwtCookie ...
func AddJwtCookie(res WebResponse, j string) {
  cok := NewCookie(WithJwtCookie(j))
  res.AddCookie(cok)
}

func DelJwtCookie(res WebResponse) {
  cok := NewCookie(WithDelJwtCookie())

  res.AddCookie(cok)
}

type (
  optionCok   func(o *optionCokTx)
  optionCokTx struct {
    name, value, path        string
    secure, httpOnly, quoted bool
    sameSite                 http.SameSite
    maxAge                   int
  }
)
