package rpc

import (
  "fmt"
  "reflect"
  "strings"
  "sync"

  "github.com/yolksys/emei/errs"
  "github.com/yolksys/emei/rpc/errors"
  "github.com/yolksys/emei/rpc/session"
)

// RegRcvr ...
func RegRcvr(v ...any) error {
  for _, value := range v {
    parse(value)
  }

  return nil
}

// ...
func Start(wg_ *sync.WaitGroup) error {
  wg_.Add(1)
  _wg = wg_

  ch, err := session.ListenSesn()
  if err != nil {
    panic(err)
  }

  go func(ch <-chan session.SesnIx) {
    for v := range ch {
      go callFunc(v)
    }
  }(ch)

  return nil
}

// callFunc ...
func callFunc(sesn session.SesnIx) {
  defer sesn.Finish()

  e := sesn.Env()
  var ci_ CallInfo
  err := sesn.HearV(&ci_)
  e.AssertErr(err, errors.ERR_ID_RPC_DECODE_CALLINFO)

  rvi := strings.Split(string(ci_), ".")
  if len(rvi) != 2 {
    e.AssertErr(errs.Wrap(fmt.Errorf("callinfo have error info: %s", ci_),
      errors.ERR_ID_RPC_CALLINFO_LEN))
  }

  rcvr, ok := _recvs[rvi[0]]
  e.AssertBool(ok, errors.ERR_ID_RPC_CALLINFO_RCVR, "have no rcvr for %s", rvi[0])
  met, ok := rcvr.funcs[rvi[1]]
  e.AssertBool(ok, errors.ERR_ID_RPC_CALLINFO_METH, "have no meth for %s", rvi[1])
  pty := rcvr.params[rvi[1]]
  lpk := pty[len(pty)-1].Kind()
  if lpk == reflect.Interface {
    pty = pty[:len(pty)-1]
  }
  params, err := sesn.Hear(pty...)
  e.AssertErr(err, errors.ERR_ID_RPC_DECODE_PARAMS)
  params = append([]reflect.Value{rcvr.value, reflect.ValueOf(e)}, params...)
  rets := met.Call(params)
  // _ = rets
}

// ...
func Shutdown() {
  // wait
  _wg.Done()
}
