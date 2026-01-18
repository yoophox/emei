package rpc

import (
  "fmt"
  "io"
  "reflect"
  "strings"
  "sync"

  "github.com/yolksys/emei/env"
  "github.com/yolksys/emei/errs"
  "github.com/yolksys/emei/kube"
  "github.com/yolksys/emei/log"
  "github.com/yolksys/emei/rpc/errors"
  "github.com/yolksys/emei/rpc/net"
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

  snet, err := kube.LookupNet("@@self")
  if err != nil {
    panic(err)
  }
  ch, err := net.Listen(":" + snet.Port)
  if err != nil {
    panic(err)
  }

  go func(ch <-chan io.ReadWriteCloser) {
    for conn := range ch {
      sess, err := newSession(conn)
      if err != nil {
        log.Event("handle new conn failed", err.Error())
        continue
      }

      dispatch(sess)
    }

    log.Event("!!!!!", "listen ch closed")
  }(ch)

  return nil
}

// dispatch ...
func dispatch(sess *session) {
  for {
    e := env.New(sess)
    go callFunc(e, sess)
  }
}

// callFunc ...
func callFunc(e env.Env, s *session) {
  fps := make([]any, 0, 16)
  action := func() error {
    return s.encode(fps)
  }

  e.Finish(action)
  defer e.Finish()

  var ci_ CallInfo
  err := s.Decode(&ci_)
  if err != nil {
    s.Close()
    e.AssertErr(errs.Wrap(fmt.Errorf("read callinfo failed"),
      errors.ERR_ID_RPC_DECODE_CALLINFO))
    return
  }

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
  params, err := s.decode(pty...)
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
