package pki

type (
  options struct {
    typ  cryptoClass
    bits int
  }
  Option func(opt *options)
)

// WithRSA ...
func WithRSA() Option {
  return func(opt *options) {
    opt.typ = _CRYPTO_CLASS_RSA
  }
}

// WithED25519 ...
func WithED25519() Option {
  return func(opt *options) {
    opt.typ = _CRYPTO_CLASS_ED25519
  }
}

// WithBits ...
func WithBits(bits int) Option {
  return func(opt *options) {
    opt.bits = bits
  }
}
