package grpc

import "github.com/yolksys/emei/rpc/call"

func init() {
  call.RegSender("grpc", newSender())
}
