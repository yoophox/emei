package svr

import "reflect"

func (l *linkTx) Values(typs ...reflect.Type) ([]reflect.Value, error) {
  ret := []reflect.Value{}

  if len(typs) == 1 && (typs[0] == typeOfReader ||
    typs[0] == typeOfWriter || typs[0] == typeOfReaderWriter) {
    ret = append(ret, reflect.ValueOf(l))
  }

  for k, v := range typs {
    vk := v.Kind()
    var vv reflect.Value
    if vk == reflect.Pointer {
      vv = reflect.New(v.Elem())
    } else {
      vv = reflect.New(v)
    }

    err := l.cc.Decode(vv)
    if err != nil {
      zero(&ret, typs[k:]...)
      return ret, err
    }

    if vk != reflect.Pointer {
      vv = vv.Elem()
    }
    ret = append(ret, vv)
  }
  return nil, nil
}

func (l *linkTx) Release() {
  l.pol.Put(l)
}
