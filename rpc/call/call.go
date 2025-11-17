package call

import (
  "fmt"
  "io"
  "reflect"

  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/log"
  "github.com/yolksys/emei/rpc/coder"
)

type Response interface {
  RValues(rTyp ...reflect.Type) ([]reflect.Value, error)
  Close()
}

type Sender interface {
  Send(e env.Env, svc, met string, args ...any) Response
  SendWithStream(e env.Env, svc, met string, args ...any) io.ReadWriteCloser
}

// RegSender ...
func RegSender(rpcName string, snd Sender) {
  _senders[rpcName] = snd
}

// call ...
func Call(e env.Env, svcName, met string, args ...any) error {
  defer e.Return()

  res := send(e, svcName, met, args...)
  defer res.Close()
  _, err := res.RValues()
  return err
}

func Call1[T1 any](e env.Env, svcName, met string, args ...any) (T1, error) {
  defer e.Return()

  t1 := reflect.TypeFor[T1]()
  res := send(e, svcName, met, args...)
  defer res.Close()

  rvs, err := res.RValues(t1)
  return rvs[0].Interface().(T1), err
}

func Call2[T1, T2 any](e env.Env, svcName, met string, args ...any) (T1, T2, error) {
  defer e.Return()

  t1 := reflect.TypeFor[T1]()
  t2 := reflect.TypeFor[T2]()
  res := send(e, svcName, met, args...)
  defer res.Close()

  rvs, err := res.RValues(t1, t2)

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

  rvs, err := res.RValues(reflect.TypeFor[T1](),
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

  rvs, err := res.RValues(reflect.TypeFor[T1](),
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

  rvs, err := res.RValues(reflect.TypeFor[T1](),
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

  rvs, err := res.RValues(reflect.TypeFor[T1](),
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

  rvs, err := res.RValues(reflect.TypeFor[T1](),
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

  rvs, err := res.RValues(reflect.TypeFor[T1](),
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

  rvs, err := res.RValues(reflect.TypeFor[T1](),
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

// CallWithRStream ...
func CallWithRStream(e env.Env, svc, met string, args ...any) (io.ReadCloser, error) {
  return callWithStream(e, svc, met, args...)
}

func CallWithWStream(e env.Env, svc, met string, args ...any) (io.WriteCloser, error) {
  return callWithStream(e, svc, met, args...)
}

func CallWithRWStream(e env.Env, svc, met string, args ...any) (io.ReadWriteCloser, error) {
  return callWithStream(e, svc, met, args...)
}

// callWithStream ...
func callWithStream(e env.Env, svc, met string, args ...any) (io.ReadWriteCloser, error) {
  defer e.Return()

  for _, v := range PriorityRpc {
    snd, ok := _senders[v]
    if !ok {
      continue
    }

    io_ := snd.SendWithStream(e, svc, met, args...)
    if e.HasError() {
      log.Debug("msg", "stream send error", "sender", v, "par", e.Err())
      continue
    }

    e.ResetErr()
    return io_, nil
  }
  return nil, fmt.Errorf("fail:no sender success, senders:%+v", _senders)
}

// send ...
func send(e env.Env, svc, met string, args ...any) Response {
  defer e.Return()

  for _, v := range PriorityRpc {
    snd, ok := _senders[v]
    if !ok {
      continue
    }

    res := snd.Send(e, svc, met, args...)
    if e.HasError() {
      e.Event("msg", "send error", "sender", v)
      continue
    }

    e.ResetErr()
    return res
  }

  e.ResetErr()
  return &defaultRes{fmt.Errorf("no sender")}
}

type defaultRes struct {
  e error
}

func (d *defaultRes) RValues(typs ...reflect.Type) ([]reflect.Value, error) {
  return coder.DefaultValues(nil, typs...), d.e
}

func (d *defaultRes) Close() {
}

var (
  _senders    = map[string]Sender{} // key=rpcName
  PriorityRpc = []string{"web"}
)
