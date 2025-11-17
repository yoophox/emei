package etcintra

import (
  "context"

  "github.com/yolksys/emei/etc/etcout"
)

type Client interface {
  GrantLease(ttl int64) Lease
  Put(ctx context.Context, typ, set, key, value string, // set can be ""
    opts ...*Option) (old string, e error) // lease as a param
  KeepLease(int64) error
  Watch(typ, set string, key ...string) (<-chan etcout.Event, error) // typ = service:cfg/service:ip, key can ve ""
  Get(typ, set string, key ...string) (map[string]string, error)     // key can be "", if that get all keyvalues of typ
}

type Lease interface {
  GetID() int64
  // TTL()
}

type Option struct {
  Typ   string
  Value any
}
