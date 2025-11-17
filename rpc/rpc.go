package rpc

import (
  "flag"

  "github.com/yolksys/emei/rpc/grpc"
  "github.com/yolksys/emei/rpc/nrpc"
  "github.com/yolksys/emei/rpc/rpcabs"
  "github.com/yolksys/emei/rpc/web"
)

type (
  RPC    = rpcabs.RPC
  Recver = rpcabs.Recver
)

func init() {
}

var (
  RpcImps   = map[string]RPC{"web": web.NewRpc(), "nrpc": nrpc.NewRpc(), "grpc": grpc.NewRpc()}
  ParseRcvr = rpcabs.ParseRcvr
  _flagSet  = flag.NewFlagSet("rpc", flag.PanicOnError)
)
