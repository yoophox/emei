package rpc

import (
  "io"
  "reflect"

  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/rpc/codec"
  "github.com/yolksys/emei/rpc/codec/gob"
  "github.com/yolksys/emei/rpc/errors"
)

type CallInfo string // rcvr.met
// newCallInfo ...
func newCallInfo(met string) string {
  return met
}

type rcvr struct {
  params map[string][]reflect.Type
  rets   map[string][]reflect.Type
  value  reflect.Value
  typ    reflect.Type
  funcs  map[string]reflect.Value
}

type codec_ interface {
  Encode(any) error
  Decode(any) error
}

type session struct {
  codec_
  io.ReadWriteCloser
  buf [1]byte
}

// newCaller ...
// read some byte from io and create codec
func newSession(io io.ReadWriteCloser) (*session, error) {
  c := &session{ReadWriteCloser: io}
  n, err := io.Read(c.buf[:])
  if err != nil {
    return nil, err
  }
  if n != 1 {
    return nil, env.Errorf(errors.ERR_ID_RPC_NEW_SESSION, "read error num datas:%d", n)
  }

  switch c.buf[0] {
  case codec.CODEC_GOB:
    c.codec_ = gob.New(io)
  case codec.CODEC_GRPC:
    break
  case codec.CODEC_HTTP_JSON:
  default:
    return nil, env.Errorf(errors.ERR_ID_RPC_NEW_SESSION, "not supported codec:%d", c.buf[0])
  }

  return c, nil
}
