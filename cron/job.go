package cron

import (
  "reflect"
  "time"
)

type job struct {
  spec   string
  ct     crontabIx
  do     reflect.Value
  params []reflect.Value
}

func (i *job) Run() *time.Time {
  s := i.do.Call(i.params)[0].Interface().(string)
  if s != "" && s != i.spec {
    i.spec = s
    var err error
    i.ct, err = parse(s)
    if err != nil {
      return nil
    }
  }
  t := time.Now()
  i.ct.Next(t)
  if t.IsZero() {
    return nil
  }
  return &t
}
