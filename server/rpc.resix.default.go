package svr

import "reflect"

func (d *defaultResIx) Values(typs ...reflect.Type) (rets []reflect.Value, err error) {
  rets = []reflect.Value{}
  for _, v := range typs {
    rets = append(rets, reflect.Zero(v))
  }
  err = d.err
  return
}

func (d *defaultResIx) Release() {
}
