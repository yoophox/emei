package coder

import (
  "fmt"
  "reflect"

  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/rpc/rpcabs"
)

// Header ...
func Header(d Decoder) (*env.Tjatse, error) {
  var h env.Tjatse
  err := d.Decode(&h)
  return &h, err
}

func Values(dec Decoder, typs ...reflect.Type) ([]reflect.Value, error) {
  rvs := make([]reflect.Value, 0)
  rl_ := len(typs)
  for i := 0; i < rl_; i++ {
    var retv reflect.Value
    retIsValue := false

    if typs[i].Kind() == reflect.Pointer {
      retv = reflect.New(typs[i].Elem())
    } else {
      retv = reflect.New(typs[i])
      retIsValue = true
    }

    err1 := dec.Decode(retv.Interface())
    if err1 != nil {
      return nil, fmt.Errorf("fail:values, at:%d, reason:%s", i, err1.Error())
    }
    if retIsValue {
      retv = retv.Elem()
    }
    rvs = append(rvs, retv)
  }

  return rvs, nil
}

func DefaultValues(_ Decoder, typs ...reflect.Type) []reflect.Value {
  rvs := make([]reflect.Value, 0, 6)
  rl_ := len(typs)
  for i := 0; i < rl_; i++ {
    var retv reflect.Value
    retIsValue := false

    if typs[i].Kind() == reflect.Pointer {
      retv = reflect.Zero(typs[i])
    } else {
      retv = reflect.New(typs[i])
      retIsValue = true
    }

    if retIsValue {
      retv = retv.Elem()
    }
    rvs = append(rvs, retv)
  }

  return rvs
}

func ParseParam(e env.Env, met *rpcabs.MethodType,
  recv *rpcabs.Recver,
) []reflect.Value {
  defer e.Return()

  v, e_ := Values(e.GetDec(), met.ArgType...)
  e.AssertErr(e_)
  vs_ := make([]reflect.Value, 0, 9)
  vs_ = append(vs_, recv.Val, reflect.ValueOf(e))
  return append(vs_, v...)
}

// AnyArrToRvlaues ...
func JsonStrArrToRvlaues(v []string, typs ...reflect.Type) ([]reflect.Value, error) {
  if len(v) != len(typs) {
    return nil, fmt.Errorf("fail:StrArrToRvalues, msg:len of v != len of typs, v:%+v, typs:%+v", v, typs)
  }
  dec := &JsonStrArrDecoder{a: v}
  return Values(dec, typs...)
}
