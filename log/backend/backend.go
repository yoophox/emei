package backend

import (
  _ "github.com/yolksys/emei/log/backend/console"
  "github.com/yolksys/emei/log/backend/intra"
  _ "github.com/yolksys/emei/log/backend/otel"
)

func Get(name string) intra.BcdNew {
  f, _ := intra.RegisteredBck[name]
  return f
}

type Backend = intra.Backend
