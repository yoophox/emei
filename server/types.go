package svr

import (
  "net"
  "net/http"
  "reflect"
  "sync"
  "time"

  "github.com/quic-go/quic-go"
)

type (
  WebSock interface {
    ReadMessage() (msgTyp int, b []byte, err error)
    ReadJSON(v any) error
    WriteMessage(msgTyp int, b []byte) error
    WriteJSON(v any) error
    // SetDeadline(t time.Time) error
    SetReadDeadline(t time.Time) error
    SetWriteDeadline(t time.Time) error
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
  // io.ReadWriteCloser
  net.Conn
  pol        *sync.Pool
  isReleased bool
}

type (
  tlsAlpnTx string
  netTx     byte
)

type quicConn struct {
  *quic.Stream
  *quic.Conn
}

type webResponsImpl struct {
  headers map[string]string
  cookies []*http.Cookie
  content [][]byte
}
