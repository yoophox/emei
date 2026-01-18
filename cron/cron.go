package cron

type Cron interface {
  Add(spec string, f CronFunc, p ...any) Canceler
}

type Canceler func()

// @return: empty string or new spec
type CronFunc func(p ...any) string

// New ...
func New() Cron {
  return nil
}
