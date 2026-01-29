package source

import (
  "fmt"

  "github.com/yoophox/emei/cfg/source/cetc"
  "github.com/yoophox/emei/cfg/source/inter"
  "github.com/yoophox/emei/cfg/source/local"
)

// Load ...
func Load(sTyp, url string) (inter.Source, error) {
  ld_, ok := _loaders[sTyp]
  if !ok {
    return nil, fmt.Errorf("fail:source->load, msg: have no loader for '%s'", sTyp)
  }
  return ld_(url)
}

type loadF func(string) (inter.Source, error)

var _loaders = map[string]loadF{
  "local": local.Load,
  "etc":   cetc.Load,
}

const (
  CFG_SOURCE_LOCAL = "local"
  CFG_SOURCE_ETC   = "etc"
  CFG_SOURCE_KUBE  = "kube"
)
