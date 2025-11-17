package web

import (
  "net/http"

  "github.com/yolksys/emei/rpc/rpcabs"
)

type webrpc struct {
  // met = http/https
  rpcabs.RpcImp
  real *http.Server
}
