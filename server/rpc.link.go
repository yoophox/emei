package svr

import (
  "fmt"
  "reflect"
  "time"

  "github.com/yoophox/emei/env"
  "github.com/yoophox/emei/errs"
)

func (l *linkTx) Values(typs ...reflect.Type) ([]reflect.Value, error) {
  ret := []reflect.Value{}

  var tja env.Tjatse
  err := l.cc.Decode(&tja)
  if err != nil {
    l.isReleased = true
    zero(&ret, typs...)
    return ret, errs.Wrap(err, ERR_ID_RPC_DECODE_TJAX)
  }
  if tja.Code != "" {
    l.isReleased = true
    zero(&ret, typs...)
    return ret, errs.Wrap(fmt.Errorf("%s", tja.Reason), errs.ErrId(tja.Code))
  }

  if len(typs) == 1 && (typs[0] == typeOfReader ||
    typs[0] == typeOfWriter || typs[0] == typeOfReaderWriter) {
    l.isReleased = true

    ret = append(ret, reflect.ValueOf(l))
    return ret, nil
  }

  for k, v := range typs {
    vv, err := decodeType(v, l.cc)
    if err != nil {
      zero(&ret, typs[k:]...)
      l.isReleased = true
      return ret, err
    }
    ret = append(ret, vv)
  }
  return nil, nil
}

func (l *linkTx) Release() {
  if l.isReleased {
    return
  }
  l.pol.Put(l)
  l.SetDeadline(time.Time{})
}
