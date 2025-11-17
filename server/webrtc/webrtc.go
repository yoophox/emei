package webrtc

import (
  "net"

  "github.com/yolksys/emei/log"
  "github.com/pion/stun"
)

// Serve ...
func Serve() {
  // Listen for UDP packets on the default STUN port (3478)
  conn, err := net.ListenPacket("udp4", ":3479")
  if err != nil {
    log.Error("fail", "listenpacket", "reason", err.Error())
  }
  defer conn.Close()

  log.Event("listenpacket at", "3479")

  buf := make([]byte, _maxPacketLen) // Buffer for incoming packets

  for {
    n, addr, err := conn.ReadFrom(buf)
    if err != nil {
      log.Event("fail", "ReadFrom", "reason", err.Error())
      continue
    }

    if stun.IsMessage(buf[:n]) {
      retv, err := stunDo(buf[:n], addr)
      if err != nil {
        log.Event("fail:stunDO, reason:%s", err.Error())
        retv, err = buildErrRetv()
      }
      n, err = conn.WriteTo(retv, addr)
      if n < len(retv) || err != nil {
        log.Event("fail:writeto, reason:%d:%s, retvl:%d", n, err.Error(), len(retv))
      }

    }
  }
}

var _maxPacketLen int = 1500
