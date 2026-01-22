package codec

import (
  "net/http"

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
