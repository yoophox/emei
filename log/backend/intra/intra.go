package intra

import (
  "context"

  "github.com/yoophox/emei/log/cache"
)

type Backend interface {
  Write(*cache.LogRecord)
}
type BcdNew func(context.Context) Backend

var RegisteredBck = map[string]BcdNew{}
