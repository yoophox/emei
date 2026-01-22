package svr

import (
  "runtime"
  "sync"

  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/kube"
)

// talk ...
func talk(e env.Env, peer string, topic string, content ...any) resIx {
  defer e.Return()
  l, err := getLink(peer)
  if err != nil {
    return &defaultResIx{err: err}
  }

  tja := env.Tjatse{
    Mid: e.TID(),
    Jwt: e.JWT().Raw(),
  }

  err = l.cc.Encode(&tja)
  if err != nil {
    return &defaultResIx{err}
  }
  err = l.cc.Encode(topic)
  if err != nil {
    return &defaultResIx{err}
  }

  return l
}

// getSesnBySvc ...
func getLink(svc string) (*linkTx, error) {
  _poolMux.RLock()
  p, ok := _linkPoolBySvc[svc]
  _poolMux.RUnlock()
  if ok {
    return poolOrNew(p, svc)
  }

  _poolMux.Lock()
  p, ok = _linkPoolBySvc[svc]
  if ok {
    _poolMux.Unlock()
    return poolOrNew(p, svc)
  }

  p = &sync.Pool{}
  _linkPoolBySvc[svc] = p
  _poolMux.Unlock()
  return poolOrNew(p, svc)
}

// poolOrNew ...
func poolOrNew(p *sync.Pool, svc string) (*linkTx, error) {
  s := p.Get()
  if s != nil {
    return s.(*linkTx), nil
  }

  si, err := kube.LookupServer(svc)
  if err != nil {
    return nil, err
  }
  conn, err := dialQuic(si.IP + ":" + si.Port)
  if err != nil {
    return nil, err
  }

  fin := func(l *linkTx) {
    l.ReadWriteCloser.Close()
  }
  cc := newGobCodec(conn)
  l := &linkTx{cc, conn, p, false}
  runtime.SetFinalizer(l, fin)
  return l, nil
}

var (
  _linkPoolBySvc = map[string]*sync.Pool{}
  _poolMux       = &sync.RWMutex{}
)
