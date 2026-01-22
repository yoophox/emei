package svr

import (
  "io"
  "net/http"
  "reflect"
  "sync"
  "time"
)

type (
  WebSock interface {
    ReadMessag() (msgTyp int, b []byte, err error)
    ReadJSON(v any) error
    WriteMessag(msgTyp int, b []byte) error
    WriteJSON(v any) error
    SetDeadTime(t time.Time) error
    SetReadDeadTime(t time.Time) error
    SetWriteDeadTime(t time.Time) error
    Close() error
  }

  SvrFor string

  WebResponse interface {
    AddHeader(k, v string)
    AddCookie(c *http.Cookie)
    Write([]byte)
  }
)

type rcvrTx struct {
  params map[string][]reflect.Type
  rets   map[string][]reflect.Type
  value  reflect.Value
  typ    reflect.Type
  funcs  map[string]reflect.Value
}

type resIx interface {
  Values(typs ...reflect.Type) ([]reflect.Value, error)
  Release()
}

type (
  defaultResIx struct{ err error }
)

type linkTx struct {
  cc codecIx
  io.ReadWriteCloser
  pol        *sync.Pool
  isStreamed bool
}

type (
  tlsAlpnTx string
  netTx     byte
)

type webResponsImpl struct{}
