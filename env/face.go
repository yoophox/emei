package env

import (
  "github.com/yoophox/emei/errs"
  "github.com/yoophox/emei/jwt"
)

type Env interface {
  //  Finish(f ...Action) // f,params...
  Trace()
  Assert()
  AssertErr(err error, eid ...errs.ErrId) // o=func() or ErrId
  AssertBool(ok bool, eid errs.ErrId, fmt_ string, args ...any)
  Done() <-chan struct{}
  IsDone() bool
  Go(f any, args ...any)
  Err() error
  TID() string
  JWT(j ...any) jwt.JWT // j=string/JWT
  Wait()                // wait for sub env
  // Clone() Env
  // Propagate(crr ...Carrier) error
}

var New = new
