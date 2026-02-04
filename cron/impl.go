package cron

import (
  "fmt"
  "reflect"
  "time"
)

// parse ...
func parse(spec string) (crontabIx, error) {
  return parser.Parse(spec)
}

type cron struct{}

func (c *cron) Add(spec string, f any, p ...any) (Canceler, error) {
  ct_, err := parse(spec)
  if err != nil {
    return nil, err
  }

  return doAdd(spec, ct_, f, p...)
}

func (c *cron) At(t time.Time, f any, p ...any) (Canceler, error) {
  return doAdd(t, true, f, p...)
}

func (c *cron) After(seconds int, f any, p ...any) (Canceler, error) {
  return doAdd(time.Now().Add(time.Duration(seconds)*time.Second), true, f, p...)
}

// doAdd ...
func doAdd(spec any, v any, f any, p ...any) (Canceler, error) {
  err := checkF(f, p...)
  if err != nil {
    return nil, err
  }

  params := []reflect.Value{}
  for _, v := range p {
    params = append(params, reflect.ValueOf(v))
  }
  j := &job{
    do:     reflect.ValueOf(f),
    params: params,
  }

  var ct_ crontabIx
  var t time.Time
  switch m := v.(type) {
  case crontabIx:
    j.spec = spec.(string)
    ct_ = m
    j.ct = m
    t = time.Now()
    t = ct_.Next(t)
  case bool:
    j.once = true
    t = spec.(time.Time)
  }

  e := AddJob(j, &t)
  return func() { _removeElemChan <- e }, nil
}

// checkF ...
func checkF(f any, args ...any) error {
  ft_ := reflect.TypeOf(f)
  if ft_.Kind() != reflect.Func {
    return fmt.Errorf("f is not function")
  }

  in_ := ft_.NumIn()
  if in_ != len(args) {
    return fmt.Errorf("parms of is not equal args:%d, %d", in_, len(args))
  }

  n := ft_.NumOut()
  if n > 1 {
    return fmt.Errorf("out num of f > 1:%d", n)
  }

  if n == 1 && ft_.Out(0).Kind() != reflect.String {
    return fmt.Errorf("The function can only return a string value")
  }

  return nil
}

var _cron = &cron{}
