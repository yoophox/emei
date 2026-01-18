package pki

type (
  options struct{}
  Option  func(opt *options)
)

// WithRSA ...
func WithRSA() Option {
  return nil
}

// WithED25519 ...
func WithED25519() Option {
  return nil
}

// WithBits ...
func WithBits(bits int) Option {
  return nil
}
