package env

import "reflect"

type Env interface {
  Finish()
  Return()
  Assert()
  HasError() bool
  ResetErr()
  Err() error
  Errorf(code uint16, f string, args ...any) error
  AssertErr(err error, clear ...func())
  AssertBool(ok bool, args ...any)
  Event(args ...interface{})
  PrintParams(v ...reflect.Value)
  GetDec() decoder
  GetMsgHeader() *Tjatse
  SetReV(v []reflect.Value)
  CheckErr(rtv ...reflect.Value) bool
}

var New = new
