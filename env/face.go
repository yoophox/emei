package env

import (
  "github.com/yolksys/emei/errs"
  "github.com/yolksys/emei/jwt"
)

type Env interface {
  //  Finish(f ...Action) // f,params...
  Return()
  Assert()
  AssertErr(err error, eid ...errs.ErrId) // o=func() or ErrId
  AssertBool(ok bool, eid errs.ErrId, fmt_ string, args ...any)
  TID() string
  JWT(j ...any) jwt.JWT // j=string/JWT
  // Propagate(crr ...Carrier) error
}

var New = new
