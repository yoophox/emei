package cron

import (
  "fmt"
  "reflect"
  "time"
)

type job struct {
  spec   string
  once   bool
  ct     crontabIx
  do     reflect.Value
  params []reflect.Value
}

func (i *job) Run() *time.Time {
  ret := i.do.Call(i.params)
  if i.once {
    return nil
  }

  if len(ret) == 1 {
    s := ret[0].Interface().(string)
    if s != "" && s != i.spec {
      i.spec = s
      var err error
      i.ct, err = parse(s)
      if err != nil {
        fmt.Println("cron.job.go parse returned spec err:" + err.Error())
        return nil
      }
    }
  }

  t := time.Now()
  t = i.ct.Next(t)
  if t.IsZero() {
    return nil
  }
  return &t
}
