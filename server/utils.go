package svr

import (
  "reflect"
  "time"
)

// timeOut ...
func rpcTimeOut() time.Time {
  return time.Now().Add(_RPC_TIMEOUT * time.Second)
}

// newType ...
func decodeType(typ reflect.Type, cc codecIx) (reflect.Value, error) {
  vk := typ.Kind()
  var vv reflect.Value
  if vk == reflect.Pointer {
    vv = reflect.New(typ.Elem())
  } else {
    vv = reflect.New(typ)
  }

  err := cc.Decode(vv)
  if err != nil {
    return vv, err
  }

  if vk != reflect.Pointer {
    vv = vv.Elem()
  }

  return vv, nil
}
