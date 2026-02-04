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
  s := fmt.Sprintf("eid:%s", e.Eid)
  if e.Err != nil {
    s += fmt.Sprintf(",by:(%s)", e.Err.Error())
  }
  return s
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

// Is ...
func Is(e error) bool {
  _, ok := e.(*Err)
  return ok
}

// Wrap ...
func Wrap(e error, id ErrId) *Err {
  if e == nil {
    return nil
  }

  return &Err{Eid: id, Err: e}
}

// ErrorF ...
func ErrorF(eid ErrId, args ...any) error {
  e := &Err{}
  e.Eid = eid
  if len(args) == 1 {
    e.Err = errors.New(args[0].(string))
  } else {
    e.Err = fmt.Errorf(args[0].(string), args[1:]...)
  }
  return e
}
