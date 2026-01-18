package env

import "github.com/yolksys/emei/errs"

type Env interface {
  Finish(f ...Action) // f,params...
  Return()
  Assert()
  AssertErr(err error) // o=func() or ErrId
  AssertBool(ok bool, eid errs.ErrId, fmt_ string, args ...any)
  Propagate(crr ...Carrier) error
}

type Carrier interface {
  Inject(*Tjatse) error
  Extract(*Tjatse) error
}

var (
  New           = new
  NewDftCarrier = newDftCarrier
)

type Action func() error
