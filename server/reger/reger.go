package reger

import (
  "reflect"

  "github.com/yolksys/emei/rpc"
  "github.com/yolksys/emei/utils"
)

// RegRpc ...
func RegRpc(t rpc.RPC) {
  RegedRpc = append(RegedRpc, t)
  SelectRpcErr = append(SelectRpcErr, reflect.SelectCase{
    Dir:  reflect.SelectRecv,
    Chan: reflect.ValueOf(t.ErrChan()),
    Send: reflect.ValueOf(nil),
  })
  err := t.Init()
  utils.AssertErr(err)
}

var (
  RegedRpc     = make([]rpc.RPC, 0, 10)
  SelectRpcErr = []reflect.SelectCase{}
)
