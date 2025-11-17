package rpcabs

import (
  "crypto/tls"
  "reflect"
  "sync"
)

type RPC interface {
  Name() string
  Start() error
  RegRcvr(rcvrs map[string]*Recver)
  Close()
  ErrChan() chan error
  IsClosed() bool
  // Recvers() map[string]*Recver
  Init() error
}

type RpcImp struct {
  // met = http/https
  Name, Port, Key, Cert, Trans string
  Wg_                          *sync.WaitGroup
  Recvs                        map[string]*Recver
  Err                          chan error // fmt: "status:,msg:"
  IsClosed_                    bool
  Certificate                  tls.Certificate
}

type Recver struct {
  Name string
  Val  reflect.Value
  Typ  reflect.Type
  Mets map[string]*MethodType
}

type MethodType struct {
  sync.Mutex // protects counters
  Method     reflect.Method
  ArgType    []reflect.Type
  ReplyType  []reflect.Type
  NumCalls   uint
}
