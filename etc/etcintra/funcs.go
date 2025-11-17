package etcintra

// Options for new client
// withPathOption ...
func WithPathOption(paths map[string]string) *Option {
  return &Option{
    Typ:   "path",
    Value: paths,
  }
}

// Options for put
// withLeaseIdOption ...
func WithLeaseIdOption(id int64) *Option {
  return &Option{
    Typ:   "leaseid",
    Value: id,
  }
}
