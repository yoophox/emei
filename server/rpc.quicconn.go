package svr

import "github.com/quic-go/quic-go"

// wrapQuicConn ...
func wrapQuicConn(s *quic.Stream, c *quic.Conn) *quicConn {
  return &quicConn{s, c}
}
