package grpc

import (
  "context"
  "net"

  "github.com/quic-go/quic-go"
)

type QuicListener struct {
  ql   *quic.Listener
  conn quic.Connection
}

// Accept waits for and returns the next connection to the listener
func (ql *QuicListener) Accept() (net.Conn, error) {
  stream, err := ql.conn.AcceptStream(context.Background())
  if err != nil {
    return nil, err
  }

  return &QuicConn{ql.conn, stream}, nil
}

// Close closes the listener
func (ql *QuicListener) Close() error {
  return nil
}

// Addr returns the listeners network address
func (ql *QuicListener) Addr() net.Addr {
  return ql.ql.Addr()
}

func Listen(ql *quic.Listener, conn quic.Connection) net.Listener {
  return &QuicListener{ql, conn}
}
