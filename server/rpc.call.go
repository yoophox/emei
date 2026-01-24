package svr

import (
  "reflect"

  "github.com/yoophox/emei/env"
)

// call ...
func Call(e env.Env, svcName, met string, args ...any) error {
  defer e.Trace()

  res := talk(e, svcName, met, args...)
  defer res.Release()
  _, err := res.Values()
  return err
}

func Call1[T1 any](e env.Env, svcName, met string, args ...any) (T1, error) {
  defer e.Trace()

  t1 := reflect.TypeFor[T1]()
  res := talk(e, svcName, met, args...)
  defer res.Release()

  rvs, err := res.Values(t1)
  return rvs[0].Interface().(T1), err
}

func Call2[T1, T2 any](e env.Env, svcName, met string, args ...any) (T1, T2, error) {
  defer e.Trace()

  t1 := reflect.TypeFor[T1]()
  t2 := reflect.TypeFor[T2]()
  res := talk(e, svcName, met, args...)
  defer res.Release()

  rvs, err := res.Values(t1, t2)

  return rvs[0].Interface().(T1),
    rvs[1].Interface().(T2),
    err
}

func Call3[T1, T2, T3 any](e env.Env, svcName, met string,
  args ...any,
) (T1, T2, T3, error) {
  defer e.Trace()

  res := talk(e, svcName, met, args...)
  defer res.Release()

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
  res := talk(e, svcName, met, args...)
  defer res.Release()

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
  res := talk(e, svcName, met, args...)
  defer res.Release()

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
  res := talk(e, svcName, met, args...)
  defer res.Release()

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
  res := talk(e, svcName, met, args...)
  defer res.Release()

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
  res := talk(e, svcName, met, args...)
  defer res.Release()

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
  res := talk(e, svcName, met, args...)
  defer res.Release()

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
