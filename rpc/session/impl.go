package session

import (
  "io"
  "reflect"
  "sync"
  "sync/atomic"

  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/rpc/codec"
)

type sesnTx struct {
  codec.CodecIx
  sync.WaitGroup
  cnt  int32
  sio  io.Closer
  env  env.Env
  iss  bool // is server?
  pool *sync.Pool
}

func (s *sesnTx) Sio() any {
  s.Add(1)
  atomic.AddInt32(&s.cnt, 1)
  return s.sio
}

func (s *sesnTx) Finish() {
  if s.cnt != 1 {
    return
  }

  s.Add(-int(s.cnt))

  if s.iss {
    return
  }

  if s.pool != nil {
    s.pool.Put(s)
  }

  // release to pool
}

func (s *sesnTx) Release() {
}

func (s *sesnTx) Close() {
  if s.sio != nil {
    s.sio.Close()
    s.Add(-int(s.cnt))
  }
}

func (s *sesnTx) Hear(typs ...reflect.Type) ([]reflect.Value, error) {
  return nil, nil
}

func (s *sesnTx) HearV(v any) error {
  return nil
}

func (s *sesnTx) Speak(w ...reflect.Value) error {
  return nil
}

func (s *sesnTx) SpeakV(v ...any) error {
  return nil
}

func (s *sesnTx) Env() env.Env {
  return s.env
}

// ...
func newSession(e env.Env, c codec.CodecIx, sio io.Closer, iss any) sesnIx {
  s := &sesnTx{
    CodecIx: c,
    sio:     sio,
    cnt:     1,
  }

  switch i_ := iss.(type) {
  case bool:
    s.iss = true
  case *sync.Pool:
    s.pool = i_
  }

  s.Add(1)
  return s
}
