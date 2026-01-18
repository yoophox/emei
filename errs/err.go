package errs

import (
  "errors"
  "fmt"
)

type Err struct {
  Eid ErrId
  Err error
}

type ErrId string

func (e *Err) Error() string {
  return fmt.Sprintf("eid:%s,by:(%s)", e.Eid, e.Err.Error())
}

func (e *Err) Isx(e_ any) bool {
  isEid := func(e *Err, eid ErrId) bool {
    if e.Eid == eid {
      return true
    }
    if _e, ok := e.Err.(*Err); ok {
      return _e.Isx(eid)
    }

    return false
  }
  switch _e := e_.(type) {
  case ErrId:
    return isEid(e, _e)
  case *Err:
    return isEid(e, _e.Eid)
  case error:
    return errors.Is(e.Err, _e)
  default:
    return false
  }
}

// Wrap ...
func Wrap(e error, id ErrId) *Err {
  if e == nil {
    return nil
  }

  return &Err{Eid: id, Err: e}
}
