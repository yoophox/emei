package cfg

import (
  "github.com/yoophox/emei/cfg/coder"
  "github.com/yoophox/emei/cfg/source"
)

const (
  CFG_SOURCE_LOCAL = source.CFG_SOURCE_LOCAL
  CFG_SOURCE_ETC   = source.CFG_SOURCE_ETC
  CFG_SOURCE_KUBE  = source.CFG_SOURCE_KUBE
)

const (
  CFG_CODER_JSON   = coder.CFG_CODER_JSON
  CFG_CODER_YAML   = coder.CFG_CODER_YAML
  CFG_CODER_STRUCT = coder.CFG_CODER_STRUCT
)
