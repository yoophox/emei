package emei

import (
  "github.com/yolksys/emei/env"
  _c "github.com/yolksys/emei/rpc/call"
  "github.com/yolksys/emei/rpc/web"
)

var (
  CallWithRStream  = _c.CallWithRStream
  CallWithWStream  = _c.CallWithWStream
  CallWithRWStream = _c.CallWithRWStream
)

func call(e Env, svcName, met string, args ...any) error {
  return _c.Call(e.(env.Env), svcName, met, args...)
}

func Call1[T1 any](e Env, svcName, met string, args ...any) (T1, error) {
  return _c.Call1[T1](e.(env.Env), svcName, met, args...)
}

func Call2[T1, T2 any](e Env, svcName, met string, args ...any) (T1, T2, error) {
  return _c.Call2[T1, T2](e.(env.Env), svcName, met, args...)
}

func Call3[T1, T2, T3 any](e Env, svcName, met string, args ...any) (T1, T2, T3, error) {
  return _c.Call3[T1, T2, T3](e.(env.Env), svcName, met, args...)
}

type (
  UpFile = web.UpFile
  DnFile = web.DnFiler
)
