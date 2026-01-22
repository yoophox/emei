package rpc

import (
  "io"
  "reflect"
  "time"
)

type CallInfo string // rcvr.met
// newCallInfo ...
func newCallInfo(met string) string {
  return met
}

type WebSock interface {
  ReadMessag() (msgTyp int, b []byte, err error)
  ReadJSON(v any) error
  WriteMessag(msgTyp int, b []byte) error
  WriteJSON(v any) error
  SetDeadTime(t time.Time) error
  SetReadDeadTime(t time.Time) error
  SetWriteDeadTime(t time.Time) error
  Close() error
}

type UpFile interface {
  Name()
  io.Reader
}

type rcvr struct {
  params map[string][]reflect.Type
  rets   map[string][]reflect.Type
  value  reflect.Value
  typ    reflect.Type
  funcs  map[string]reflect.Value
}
