package env

import (
  "github.com/yoophox/emei/errs"
  "github.com/yoophox/emei/jwt"
  "github.com/yoophox/emei/log"
)

type Env interface {
  //  Finish(f ...Action) // f,params...
  log.Logger
  Trace(name string)
  Assert()
  AssertErr(err error, eid ...errs.ErrId) // o=func() or ErrId
  AssertBool(ok bool, eid errs.ErrId, args ...any)
  Done() <-chan struct{}
  IsDone() bool
  Go(f any, args ...any)
  Err() error
  TID() string
  JWT(j ...any) jwt.JWT // j=string/JWT
  Wait()                // wait for sub env
  WaitAny()
  Cancel()
}

var New = new
