package nrpc

import (
  "net/rpc"
  "sync"

  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/rpc/rpcabs"
  "github.com/quic-go/quic-go"
)

type client struct {
  usedPool    *sync.Pool
  callPool    *sync.Pool
  stmCallPool *sync.Pool
  c           quic.Connection
}

type nrpcCaller struct {
  quic.Stream
  codec rpc.ClientCodec
  isStm bool
  stm   stream
  rpc.Request
  rpc.Response
  req, retv NRpcData
}

type sender struct {
  rpcabs.SndImp
}

type response struct {
  clr *nrpcCaller
}

// for request and responsexs
type NRpcData struct {
  MsgH  *env.Tjatse
  IsStm *bool
  Data  []string
}

type nrpc struct {
  rpcabs.RpcImp
}

type server struct {
  quic.Stream
  *gobServerCodec
  rpc.Response
  rpc.Request
  req, res   NRpcData
  stm        stream
  isStmResed bool
}
