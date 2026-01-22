package svr

import (
  "go/token"
  "reflect"
  "strings"

  "github.com/yolksys/emei/log"
)

func parseWeb(rcvr any) error {
  rv_ := reflect.ValueOf(rcvr)
  rt_ := reflect.TypeOf(rcvr)
  rn_ := reflect.Indirect(rv_).Type().Name()

  if rn_ == "" {
    panic("web:type name is empty")
  }
  if !token.IsExported(rn_) {
    panic("web: type is not exported")
  }

  cvx := &rcvrTx{
    params: make(map[string][]reflect.Type),
    funcs:  make(map[string]reflect.Value),
    value:  rv_,
    typ:    rt_,
  }

  mn_ := rt_.NumMethod()
  for i := range mn_ {
    m := rt_.Method(i)
    mt_ := m.Type
    mn_ := strings.ToLower(m.Name)
    if mt_.NumOut() > 0 {
      log.Debug("met", mn_, "have return value", mt_.NumOut())
      continue
    }

    if mt_.NumIn() < 3 || mt_.NumIn() > 4 {
      log.Debug("met", mn_, "numin", mt_.NumIn())
      continue
    }

    if mt_.In(1) != typeOfEnv {
      log.Debug("met", mn_, "first parma must be env", "")
      continue
    }

    if mt_.NumIn() == 3 && mt_.In(2) != typeOfWebsock {
      log.Debug("met", mn_, "three", "websock")
      continue
    }

    if mt_.NumIn() == 4 && (mt_.In(2) != typeOfWebRes || mt_.In(3) != typeOfRequest) {
      log.Debug("met", mn_, "four", "webrespones and *http.request")
      continue
    }

    cvx.funcs[mn_] = m.Func
    ps_ := []reflect.Type{}
    for i := 2; i < mt_.NumIn(); i++ {
      ps_ = append(ps_, mt_.In(i))
    }
    cvx.params[mn_] = ps_
  }

  if len(cvx.funcs) <= 0 {
    panic("web rcvr have no needed func")
  }

  log.Debug("rcvr", cvx)
  _webRecvs[strings.ToLower(rn_)] = cvx

  return nil
}
