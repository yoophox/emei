package source

import (
  "errors"
  "fmt"

  "github.com/yolksys/emei/cfg/source/cetc"
  "github.com/yolksys/emei/cfg/source/inter"
  "github.com/yolksys/emei/cfg/source/local"
)

// Load ...
func Load(sTyp, url string) (inter.Source, error) {
  ld_, ok := _loaders[sTyp]
  if !ok {
    return nil, errors.New(fmt.Sprintf("fail:source->load, msg: have no loader for '%s'", sTyp))
  }
  return ld_(url)
}

type loadF func(string) (inter.Source, error)

var _loaders = map[string]loadF{
  "local": local.Load,
  "etc":   cetc.Load,
}
