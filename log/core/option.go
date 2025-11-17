package core

type optfunc func(*option)

type option struct {
  isCacheMod bool
}

// resetOpt ...
func resetOpt(o *option) {
  o.isCacheMod = false
}

// WithCacheMode ...
func WithCacheMode() optfunc {
  return func(o *option) {
    o.isCacheMod = true
  }
}
