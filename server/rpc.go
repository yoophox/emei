package svr

import (
  "reflect"
  "strings"
  "time"

  "github.com/yoophox/emei/env"
  "github.com/yoophox/emei/errs"
  "github.com/yoophox/emei/log"
)

// dispatch ...
func dispatchRpc(l *linkTx) {
  defer func() {
    log.Event("link end", "")
  }()

  for {
    var tja env.Tjatse
    err := l.cc.Decode(&tja)
    if err != nil {
      log.Event()
      l.Conn.Close()
      return
    }
    var topic string
    l.SetReadDeadline(rpcTimeOut())
    err = l.cc.Decode(&topic)
    if err != nil {
      log.Event()
      l.Conn.Close()
      return
    }

    l.SetReadDeadline(time.Time{})
    e := env.New(&tja)
    assist(e, l, topic)
    e.Release()
  }
}

// ...
func assist(e env.Env, l *linkTx, topic string) {
  isStream := false
  var ret []reflect.Value

  defer func() {
    if isStream {
      return
    }

    var tja env.Tjatse
    tja.Mid = e.TID()
    tja.Jwt = e.JWT().Raw()
    err := e.Err()
    if err != nil {
      e_, ok := err.(*errs.Err)
      if ok {
        tja.Code = string(e_.Eid)
      }
      tja.Reason = err.Error()
    }

    l.SetDeadline(rpcTimeOut())
    e_ := l.cc.Encode(&tja)
    if e_ != nil || err != nil {
      l.Close()
      return
    }

    for _, v := range ret {
      err := l.cc.Encode(v.Interface())
      if err != nil {
        l.Close()
        return
      }
    }

    l.SetDeadline(time.Time{})
  }()
  defer e.Trace()

  pos := strings.Index(topic, ".")
  e.AssertBool(pos > 0, ERR_ID_RPC_CALLINFO_LEN, topic)
  rcvr, ok := _rpcRecvs[topic[:pos]]
  e.AssertBool(ok, ERR_ID_RPC_CALLINFO_RCVR, "no rcvr for topic: %s", topic)
  m, ok := rcvr.funcs[topic[pos+1:]]
  e.AssertBool(ok, ERR_ID_RPC_CALLINFO_RCVR, "no rcvr for func: %s", topic)
  ptyps := rcvr.params[topic[pos+1:]]
  lp_ := len(ptyps)
  if lp_ > 0 && (ptyps[lp_-1] == typeOfReader ||
    ptyps[lp_-1] == typeOfWriter || ptyps[lp_-1] == typeOfReaderWriter) {
    isStream = true
    ptyps = ptyps[:lp_-1]
  }

  l.SetDeadline(rpcTimeOut())
  params := []reflect.Value{rcvr.value, reflect.ValueOf(e)}
  if len(ptyps) > 0 {
    for _, t := range ptyps {
      vv, err := decodeType(t, l.cc)
      e.AssertErr(err, ERR_ID_RPC_DECODE_PARAMS)
      params = append(params, vv)
    }
  }

  if isStream {
    var tja env.Tjatse
    tja.Mid = e.TID()
    tja.Jwt = e.JWT().Raw()
    err := l.cc.Encode(&tja)
    e.AssertErr(err, ERR_ID_RPC_ENCODE)
  }
  ret = m.Call(params)
  if isStream {
    l.Close()
    return
  }
  e.Assert()
  if len(ret) > 0 && ret[len(ret)-1].Type() == typeOfError {
    e.AssertErr(ret[len(ret)-1].Interface().(error), ERR_ID_RPC_CALL_ERROR)
  }
}
