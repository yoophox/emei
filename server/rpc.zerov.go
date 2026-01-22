package svr

import (
  "reflect"
)

// zero ...
func zero(ret *[]reflect.Value, typs ...reflect.Type) {
  for _, v := range typs {
    *ret = append(*ret, reflect.Zero(v))
  }
}
