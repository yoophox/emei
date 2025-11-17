package nrpc

import (
  "errors"
  "reflect"
  "strings"
  "sync/atomic"

  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/rpc/call"
  "github.com/yolksys/emei/rpc/coder"
)

func (s *server) serve(cnt *int32) (err error) {
  err = s.gobServerCodec.ReadRequestHeader(&s.Request)
  atomic.AddInt32(cnt, -1)
  if err != nil {
    return
  }
  s.isStmResed = false
  s.Response.Seq = s.Request.Seq
  s.res.Data = s.res.Data[:0]
  s.res.MsgH = nil
  s.req.IsStm = nil
  if s.req.MsgH != nil {
    s.req.MsgH.Jwt = ""
  }
  defer func() {
    err_ := s.served()
    if err_ != nil {
      err = err_
    }
  }()

  e := env.New("nrpc", s.Request.ServiceMethod, newEnc(s), newDec(s))
  defer e.Finish()
  e.Assert()

  _path := strings.Split(s.Request.ServiceMethod, ".")
  if len(_path) != 2 {
    e.AssertErr(errors.New("fail:met,met:" + s.Request.ServiceMethod))
  }

  recv, ok := _nrpc.Recvs[_path[0]]
  if !ok {
    e.AssertErr(errors.New("fail: no recver,name:" + _path[0]))
  }

  met, ok := recv.Mets[_path[1]]
  if !ok {
    e.AssertErr(errors.New("fail: no method,name:" + _path[1]))
  }
  parms := coder.ParseParam(e, met, recv)
  e.Assert()
  e.PrintParams(parms[2:]...)
  if s.req.IsStm != nil && *s.req.IsStm {
    err := s.stmResonse()
    e.AssertErr(err)
    s.isStmResed = true
    parms = append(parms, reflect.ValueOf(&s.stm))
  }

  f := met.Method.Func
  rtv := f.Call(parms)
  if s.req.IsStm == nil || !*s.req.IsStm {
    e.SetReV(rtv)
    return nil
  }
  if e.CheckErr(rtv...) {
    err = s.stm.writeErr(e, e.Err())
    e.ResetErr()
  }
  return
}

func init() {
  call.RegSender("nrpc", newSender())
}
