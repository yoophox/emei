package cron

import "time"

type Cron interface {
  Add(spec string, f any, p ...any) (Canceler, error)
  At(t time.Time, f any, param ...any) (Canceler, error)
  After(seconds int, f any, param ...any) (Canceler, error)
}

type Canceler func()

// @return: empty string or new spec
type CronFunc func(p ...any) string

// New ...
func New() Cron {
  return _cron
}
