package coder

import (
  "errors"
  "fmt"

  "github.com/yolksys/emei/cfg/coder/cfgc"
  "github.com/yolksys/emei/cfg/coder/ckube"
  "github.com/yolksys/emei/cfg/coder/yaml"
  "github.com/yolksys/emei/cfg/source/inter"
  "github.com/yolksys/emei/cfg/values"
)

// encoder ...
func Encode(eType string, s inter.Source) (values.Values, error) {
  enc, ok := _encs[eType]
  if !ok {
    return nil, errors.New(
      fmt.Sprintf("fail: code->encode, msg:cann'// TODO: find encoder for '%s'", eType))
  }

  return enc(s)
}

type encF func(s inter.Source) (values.Values, error)

var (
  _encs = map[string]encF{
    "json": cfgc.Encode,
    "yaml": yaml.Encode,
    "kube": ckube.Encode,
  }
  p = 0
)

const (
  CFG_CODER_JSON   = "json"
  CFG_CODER_YAML   = "yaml"
  CFG_CODER_STRUCT = "struct"
)
