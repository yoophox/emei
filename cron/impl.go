package cron

import (
  "reflect"
  "time"
)

// parse ...
func parse(spec string) (crontabIx, error) {
  return parser.Parse(spec)
}

type cron struct{}

func (c *cron) Add(spec string, f CronFunc, p ...any) (Canceler, error) {
  ct_, err := parse(spec)
  if err != nil {
    return nil, err
  }

  params := []reflect.Value{}
  for _, v := range p {
    params = append(params, reflect.ValueOf(v))
  }
  j := &job{
    spec:   spec,
    ct:     ct_,
    do:     reflect.ValueOf(f),
    params: params,
  }

  t := time.Now()
  ct_.Next(t)
  e := AddJob(j, &t)

  return func() { _removeElemChan <- e }, nil
}

var _cron = &cron{}
