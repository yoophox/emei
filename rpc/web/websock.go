package web

import (
  "net/http"
  "time"

  "github.com/gorilla/websocket"
)

// ...
func newWebSock(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
  upgrader := websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
  }
  return upgrader.Upgrade(w, r, nil)
}

type WebSock interface {
  ReadMessag() (msgTyp int, b []byte, err error)
  ReadJSON(v interface{}) error
  WriteMessag(msgTyp int, b []byte) error
  WriteJSON(v interface{}) error
  SetDeadTime(t time.Time) error
  SetReadDeadTime(t time.Time) error
  SetWriteDeadTime(t time.Time) error
  Close() error
}
