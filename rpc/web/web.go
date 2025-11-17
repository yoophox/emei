package web

import "github.com/yolksys/emei/rpc/call"

func init() {
  call.RegSender("web", newWebSender())
}
