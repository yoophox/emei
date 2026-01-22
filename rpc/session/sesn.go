package session

import (
  "reflect"
  "sync"

  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/kube"
  "github.com/yolksys/emei/rpc/codec"
)

type SesnIx interface {
  // Decode(v any) error
  // Encode(v any) error
  // DecodeTyps(typs ...reflect.Type) ([]reflect.Value, error)
  // EncodeValues(vs ...reflect.Value) error
  Hear(sentence ...reflect.Type) (w []reflect.Value, err error)
  HearV(v any) error
  Speak(w ...reflect.Value) error
  SpeakV(v ...any) error
  Sio() any // reader,writer, webscok,webtrans, []file
  Env() env.Env
  Wait() // for server
  Finish()
  Close()
  // Propagate()
  // Topic()
  // close()
}

// ...
func ListenSesn( /*net netTx, addr string*/ ) (<-chan SesnIx, error) {
  return nil, nil
}

// dial ...
func DialSesn(svc string, e env.Env) (SesnIx, error) {
  return nil, nil
}

// getSesnBySvc ...
func getSesnBySvc(e env.Env, svc string) (SesnIx, error) {
  _poolMux.RLock()
  p, ok := _dialedSesnPool[svc]
  _poolMux.RUnlock()
  if ok {
    return poolOrNew(e, p, svc)
  }

  _poolMux.Lock()
  p, ok = _dialedSesnPool[svc]
  if ok {
    _poolMux.Unlock()
    return poolOrNew(e, p, svc)
  }

  p = &sync.Pool{}
  _dialedSesnPool[svc] = p
  return poolOrNew(e, p, svc)
}

// poolOrNew ...
func poolOrNew(e env.Env, p *sync.Pool, svc string) (SesnIx, error) {
  s := p.Get()
  if s != nil {
    return s.(SesnIx), nil
  }

  si, err := kube.LookupServer(svc)
  if err != nil {
    return nil, err
  }
  conn, err := dialQuic(si.IP + ":" + si.Port)
  if err != nil {
    return nil, err
  }
  cc := codec.NewGob(conn)
  return newSession(e, cc, conn, p), nil
}

type (
  tlsAlpnTx string
  netTx     byte
)

const (
  RPC_NET_QUIC netTx = iota
  RPC_NET_TCP
)

const (
  RPC_CODEC_GOB        tlsAlpnTx = "gob"
  RPC_CODEC_GRPC       tlsAlpnTx = "grpc"
  RPC_CODEC_JSON       tlsAlpnTx = "json"
  RPC_CODEC_HTTP_JSON  tlsAlpnTx = "h2"
  RPC_CODEC_HTTP3_JSON tlsAlpnTx = "h3"
)

var (
  _dialedSesnPool = map[string]*sync.Pool{}
  _poolMux        = &sync.RWMutex{}
)
