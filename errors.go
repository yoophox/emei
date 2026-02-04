package emei

import "github.com/yoophox/emei/errs"

type ErrId = errs.ErrId

var (
  // Is(e error)
  IsErr = errs.Is
  // ErrorF(eid errs.ErrId, args ...any)
  Errorf = errs.ErrorF
  // Wrap(e error, id errs.ErrId)
  WrapErr = errs.Wrap
)
