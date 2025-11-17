package etcd

import (
  "context"
  "fmt"
  "reflect"
  "time"

  "go.etcd.io/etcd/client/v3"
  "go.etcd.io/etcd/client/v3/namespace"
  "github.com/yolksys/emei/etc/etcintra"
  "github.com/yolksys/emei/etc/etcout"
)

type client struct {
  *clientv3.Client
  prefixClis     map[string]*prefixedCli
  addKeepedLease chan (<-chan *clientv3.LeaseKeepAliveResponse)
  addWatch       chan *watcher
  opts           []*etcintra.Option
}

// use ...
func Use(addrs []string, opts ...*etcintra.Option) etcintra.Client {
  c, err := clientv3.New(clientv3.Config{
    Endpoints:   addrs,
    DialTimeout: 2 * time.Second,
  })
  if err != nil {
    panic(err.Error())
  }

  _clt = &client{
    Client:         c,
    opts:           opts,
    addKeepedLease: make(chan (<-chan *clientv3.LeaseKeepAliveResponse), 50),
    addWatch:       make(chan *watcher, 50),
    prefixClis:     make(map[string]*prefixedCli),
  }

  _clt.parseOpts()
  go _clt.doKeepLease()
  go _clt.doWatch()

  return nil
}

func (c *client) GrantLease(ttl int64) (etcintra.Lease, error) {
  l, err := c.Grant(context.Background(), ttl)
  return (*lease)(l), err
}

func (c *client) Put(ctx context.Context, typ, set, key, value string,
  opts ...*etcintra.Option,
) (old string, err error) {
  o := []clientv3.OpOption{}
  for _, value := range opts {
    switch value.Typ {
    case "leaseid":
      o = append(o, clientv3.WithLease(value.Value.(clientv3.LeaseID)))
    }
  }

  cli, ok := c.prefixClis[typ]
  if !ok {
    return "", fmt.Errorf("illegal type: %s", typ)
  }

  k := key
  if k != "" && set != "" {
    k = set + "/" + k
  } else {
    k = set + k
  }

  v, err := cli.Put(ctx, k, value, o...)
  if err != nil {
    return "", err
  }

  return v.PrevKv.String(), nil
}

func (c *client) KeepLease(id int64) error {
  r, err := c.KeepAlive(context.Background(), clientv3.LeaseID(id))
  if err != nil {
    return err
  }

  c.addKeepedLease <- r

  return nil
}

func (c *client) Watch(typ, set string, key ...string) (<-chan etcout.Event, error) {
  cli, ok := c.prefixClis[typ]
  if !ok {
    return nil, fmt.Errorf("fail: watch,msg: illegal type: %s", typ)
  }

  var o clientv3.OpOption = nil
  k := ""
  if key == nil {
    o = clientv3.WithPrefix()
  } else {
    k = key[0]
  }
  if k != "" && set != "" {
    k = set + "/" + k
  } else {
    k = set + k
  }

  ch_ := make(chan etcout.Event, 50)
  c.addWatch <- &watcher{
    out: ch_,
    inc: cli.Watch(context.Background(), k, o),
  }

  return ch_, nil
}

func (c *client) Get(typ, set string, key ...string) (map[string]string, error) {
  cli, ok := c.prefixClis[typ]
  if !ok {
    return nil, fmt.Errorf("fail: Get,msg: illegal type: %s", typ)
  }

  var o clientv3.OpOption = nil
  k := ""
  if key == nil {
    o = clientv3.WithPrefix()
  } else {
    k = key[0]
  }
  if k != "" && set != "" {
    k = set + "/" + k
  } else {
    k = set + k
  }

  v, err := cli.Get(context.Background(), k, o)
  if err != nil {
    return nil, err
  }

  ret := map[string]string{}

  for _, value := range v.Kvs {
    ret[string(value.Key)] = string(value.Value)
  }

  return ret, nil
}

func (c *client) doKeepLease() {
  cases := make([]reflect.SelectCase, 0, 50)
  cases = append(cases, reflect.SelectCase{
    Dir:  reflect.SelectRecv,
    Chan: reflect.ValueOf(c.addKeepedLease),
  })

  for {
    i, v, ok := reflect.Select(cases)
    if !ok {
      continue
    }

    if i != 0 {
      continue
    }

    cases = append(cases, reflect.SelectCase{
      Dir:  reflect.SelectRecv,
      Chan: v,
    })
  }
}

func (c *client) doWatch() {
  var (
    watchers = make([]*watcher, 0, 50)
    cases    = make([]reflect.SelectCase, 0, 50)
  )
  cases = append(cases, reflect.SelectCase{
    Dir:  reflect.SelectRecv,
    Chan: reflect.ValueOf(c.addWatch),
  })

  for {
    i, v, ok := reflect.Select(cases)
    if !ok {
      continue
    }

    if i == 0 {
      // do addWatch
      watchers = append(watchers, v.Interface().(*watcher))
      cases = append(cases, reflect.SelectCase{
        Dir:  reflect.SelectRecv,
        Chan: reflect.ValueOf(v.Interface().(*watcher).inc),
      })
      continue
    }

    out := watchers[i-1].out
    for _, value := range v.Interface().(*clientv3.WatchResponse).Events {
      out <- (*event)(value)
    }
  }
}

func (c *client) parseOpts() {
  for _, value := range c.opts {
    switch value.Typ {
    case "path":
      v := value.Value.(map[string]string)
      for t, p := range v {
        c.prefixClis[t] = &prefixedCli{
          KV:      namespace.NewKV(c.Client, p),
          Lease:   namespace.NewLease(c.Client, p),
          Watcher: namespace.NewWatcher(c.Client, p),
        }
      }
    }
  }
}

type watcher struct {
  out chan<- etcout.Event
  inc clientv3.WatchChan
}
type prefixedCli struct {
  clientv3.KV
  clientv3.Lease
  clientv3.Watcher
}

type (
  event clientv3.Event
  lease clientv3.LeaseGrantResponse
)

func (e *event) EType() string {
  switch e.Type {
  case clientv3.EventTypePut:
    if e.Kv.Version == 0 {
      return "new"
    }
    return "update"
  case clientv3.EventTypeDelete:
    return "delete"
  }

  return ""
}

func (e *event) EKey() string {
  return string(e.Kv.Key)
}

func (e *event) EValue() string {
  return string(e.Kv.Value)
}

func (e *event) EVersion() int64 {
  return e.Kv.Version
}

func (l *lease) GetID() int64 {
  return int64(l.ID)
}

func init() {
  etcintra.EtcBcks["etcd"] = Use
}

var (
  _clt *client
  k    int
)
