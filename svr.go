package emei

import (
  "github.com/yoophox/emei/env"
  svr "github.com/yoophox/emei/server"
)

// server
var (
  Serve    = svr.Serve
  ServeFor = svr.ServeFor
)

// call ...
func Call(e env.Env, svc, met string, args ...any) error {
  return svr.Call(e, svc, met, args...)
}

// Call1 ...
func Call1[T any](e env.Env, svc, met string, args ...any) (T, error) {
  return svr.Call1[T](e, svc, met, args...)
}

func Call2[T1, T2 any](e env.Env, svc, met string, args ...any) (T1, T2, error) {
  return svr.Call2[T1, T2](e, svc, met, args...)
}

func Call3[T1, T2, T3 any](e env.Env, svc, met string, args ...any) (T1, T2, T3, error) {
  return svr.Call3[T1, T2, T3](e, svc, met, args...)
}

func Call4[T1, T2, T3, T4 any](e env.Env, svc, met string, args ...any) (T1, T2, T3, T4, error) {
  return svr.Call4[T1, T2, T3, T4](e, svc, met, args...)
}

func Call5[T1, T2, T3, T4, T5 any](e env.Env, svc, met string, args ...any) (T1, T2, T3, T4, T5, error) {
  return svr.Call5[T1, T2, T3, T4, T5](e, svc, met, args...)
}

func Call6[T1, T2, T3, T4, T5, T6 any](e env.Env, svc, met string, args ...any) (T1, T2, T3, T4, T5, T6, error) {
  return svr.Call6[T1, T2, T3, T4, T5, T6](e, svc, met, args...)
}

func Call7[T1, T2, T3, T4, T5, T6, T7 any](e env.Env, svc, met string, args ...any) (T1, T2, T3, T4, T5, T6, T7, error) {
  return svr.Call7[T1, T2, T3, T4, T5, T6, T7](e, svc, met, args...)
}

func Call8[T1, T2, T3, T4, T5, T6, T7, T8 any](e env.Env, svc, met string, args ...any) (T1, T2, T3, T4, T5, T6, T7, T8, error) {
  return svr.Call8[T1, T2, T3, T4, T5, T6, T7, T8](e, svc, met, args...)
}

func Call9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](e env.Env, svc, met string, args ...any) (T1, T2, T3, T4, T5, T6, T7, T8, T9, error) {
  return svr.Call9[T1, T2, T3, T4, T5, T6, T7, T8, T9](e, svc, met, args...)
}
