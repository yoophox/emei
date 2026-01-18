package rpc

import (
  "reflect"
  "sync"

  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/errs"
)

// call ...
func Call(e env.Env, svcName, met string, args ...any) error {
  defer e.Return()

  res := send(e, svcName, met, args...)
  defer res.Close()
  _, err := res.Values()
  return err
}

func Call1[T1 any](e env.Env, svcName, met string, args ...any) (T1, error) {
  defer e.Return()

  t1 := reflect.TypeFor[T1]()
  res := send(e, svcName, met, args...)
  defer res.Close()

  rvs, err := res.Values(t1)
  return rvs[0].Interface().(T1), err
}

func Call2[T1, T2 any](e env.Env, svcName, met string, args ...any) (T1, T2, error) {
  defer e.Return()

  t1 := reflect.TypeFor[T1]()
  t2 := reflect.TypeFor[T2]()
  res := send(e, svcName, met, args...)
  defer res.Close()

  rvs, err := res.Values(t1, t2)

  return rvs[0].Interface().(T1),
    rvs[1].Interface().(T2),
    err
}

func Call3[T1, T2, T3 any](e env.Env, svcName, met string,
  args ...any,
) (T1, T2, T3, error) {
  defer e.Return()

  res := send(e, svcName, met, args...)
  defer res.Close()

  rvs, err := res.Values(reflect.TypeFor[T1](),
    reflect.TypeFor[T2](),
    reflect.TypeFor[T3]())

  return rvs[0].Interface().(T1),
    rvs[1].Interface().(T2),
    rvs[2].Interface().(T3),
    err
}

func Call4[T1, T2, T3, T4 any](e env.Env, svcName, met string,
  args ...any,
) (T1, T2, T3, T4, error) {
  res := send(e, svcName, met, args...)
  defer res.Close()

  rvs, err := res.Values(reflect.TypeFor[T1](),
    reflect.TypeFor[T2](),
    reflect.TypeFor[T3](),
    reflect.TypeFor[T4]())

  return rvs[0].Interface().(T1),
    rvs[1].Interface().(T2),
    rvs[2].Interface().(T3),
    rvs[3].Interface().(T4),
    err
}

func Call5[T1, T2, T3, T4, T5 any](e env.Env, svcName, met string,
  args ...any,
) (T1, T2, T3, T4, T5, error) {
  res := send(e, svcName, met, args...)
  defer res.Close()

  rvs, err := res.Values(reflect.TypeFor[T1](),
    reflect.TypeFor[T2](),
    reflect.TypeFor[T3](),
    reflect.TypeFor[T4](),
    reflect.TypeFor[T5]())

  return rvs[0].Interface().(T1),
    rvs[1].Interface().(T2),
    rvs[2].Interface().(T3),
    rvs[3].Interface().(T4),
    rvs[4].Interface().(T5),
    err
}

func Call6[T1, T2, T3, T4, T5, T6 any](e env.Env, svcName, met string,
  args ...any,
) (T1, T2, T3, T4, T5, T6, error) {
  res := send(e, svcName, met, args...)
  defer res.Close()

  rvs, err := res.Values(reflect.TypeFor[T1](),
    reflect.TypeFor[T2](),
    reflect.TypeFor[T3](),
    reflect.TypeFor[T4](),
    reflect.TypeFor[T5](),
    reflect.TypeFor[T6]())

  return rvs[0].Interface().(T1),
    rvs[1].Interface().(T2),
    rvs[2].Interface().(T3),
    rvs[3].Interface().(T4),
    rvs[4].Interface().(T5),
    rvs[5].Interface().(T6),
    err
}

func Call7[T1, T2, T3, T4, T5, T6, T7 any](e env.Env, svcName, met string,
  args ...any,
) (T1, T2, T3, T4, T5, T6, T7, error) {
  res := send(e, svcName, met, args...)
  defer res.Close()

  rvs, err := res.Values(reflect.TypeFor[T1](),
    reflect.TypeFor[T2](),
    reflect.TypeFor[T3](),
    reflect.TypeFor[T4](),
    reflect.TypeFor[T5](),
    reflect.TypeFor[T6](),
    reflect.TypeFor[T7]())

  return rvs[0].Interface().(T1),
    rvs[1].Interface().(T2),
    rvs[2].Interface().(T3),
    rvs[3].Interface().(T4),
    rvs[4].Interface().(T5),
    rvs[5].Interface().(T6),
    rvs[6].Interface().(T7),
    err
}

func Call8[T1, T2, T3, T4, T5, T6, T7, T8 any](e env.Env, svcName, met string,
  args ...any,
) (T1, T2, T3, T4, T5, T6, T7, T8, error) {
  res := send(e, svcName, met, args...)
  defer res.Close()

  rvs, err := res.Values(reflect.TypeFor[T1](),
    reflect.TypeFor[T2](),
    reflect.TypeFor[T3](),
    reflect.TypeFor[T4](),
    reflect.TypeFor[T5](),
    reflect.TypeFor[T6](),
    reflect.TypeFor[T7](),
    reflect.TypeFor[T8]())

  return rvs[0].Interface().(T1),
    rvs[1].Interface().(T2),
    rvs[2].Interface().(T3),
    rvs[3].Interface().(T4),
    rvs[4].Interface().(T5),
    rvs[5].Interface().(T6),
    rvs[6].Interface().(T7),
    rvs[7].Interface().(T8),
    err
}

func Call9[T1, T2, T3, T4, T5, T6, T7, T8, T9 any](e env.Env, svcName, met string,
  args ...any,
) (T1, T2, T3, T4, T5, T6, T7, T8, T9, error) {
  res := send(e, svcName, met, args...)
  defer res.Close()

  rvs, err := res.Values(reflect.TypeFor[T1](),
    reflect.TypeFor[T2](),
    reflect.TypeFor[T3](),
    reflect.TypeFor[T4](),
    reflect.TypeFor[T5](),
    reflect.TypeFor[T6](),
    reflect.TypeFor[T7](),
    reflect.TypeFor[T8](),
    reflect.TypeFor[T9]())

  return rvs[0].Interface().(T1),
    rvs[1].Interface().(T2),
    rvs[2].Interface().(T3),
    rvs[3].Interface().(T4),
    rvs[4].Interface().(T5),
    rvs[5].Interface().(T6),
    rvs[6].Interface().(T7),
    rvs[7].Interface().(T8),
    rvs[8].Interface().(T9),
    err
}

// ...
func send(e env.Env, svc, met string, args ...any) resIx {
  defer e.Return()

  sess, err := getSession(svc)
  e.AssertErr(err)
  err = sess.encode(newCallInfo(met))
  e.AssertErr(err)
  err = sess.encode(args...)
  e.AssertErr(err)
  return (*sessionResIx)(sess)
}

type resIx interface {
  Values(typs ...reflect.Type) ([]reflect.Value, error)
  Close()
}

type (
  defaultResIx struct{ err error }
  sessionResIx session
)

func (d *defaultResIx) Values(typs ...reflect.Type) (rets []reflect.Value, err error) {
  rets = []reflect.Value{}
  for _, v := range typs {
    rets = append(rets, reflect.Zero(v))
  }
  err = d.err
  return
}

func (d *defaultResIx) Close() {
}

func (r *sessionResIx) Values(typs ...reflect.Type) (rets []reflect.Value, err error) {
  rets, err = (*session)(r).decode(typs...)
  return
}

func (r *sessionResIx) Close() {
}

var (
  _sessionsForService = map[string]*session{}
  _sessionsMux        = sync.RWMutex{}
)

// getSession ...
func getSession(svc string) (s *session, err error) {
  return
}

// ...
func delSession(svc string) {
}
