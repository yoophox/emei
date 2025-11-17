package env

import (
  "context"
  "fmt"
  "path"
  "reflect"
  "sync"

  "github.com/yolksys/emei/log"
  "github.com/yolksys/emei/log/core"

  "github.com/yolksys/emei/otel"
  "github.com/yolksys/emei/utils"
)

type env struct {
  L   log.Logger
  Cf_ []func()        // callback
  Ecf []func()        // error call back
  ReV []reflect.Value // return values
  Par *env            // up env

  msgH  *Tjatse
  enc   encoder // encode return value
  dec   decoder //
  span  otel.Span
  rpc   string
  met   string
  meter otel.Meter
  err   error
}

// New ...
func new(rpc, met string, enc encoder, dec decoder) Env {
  e := pool.Get().(*env)
  e.L = log.New(context.Background(), core.WithCacheMode())
  e.L.CallerSkip(envLogSkip)
  e.L.Event("*", "start")
  e.enc = enc
  e.dec = dec

  if dec == nil {
    e.msgH = &Tjatse{}
  } else {
    h, err := dec.Header()
    e.msgH = h
    e.err = err
  }

  e.span = otel.Trace(e.msgH, rpc+"."+met)
  if e.span == nil {
    return e
  }
  e.L.SetTraceId(e.span.TID())
  e.L.SetTraceId(e.span.TID())
  if e.msgH.Mid == "" {
    e.span.AddAttri("uid", e.msgH.Uid())
    e.span.AddAttri("uname", e.msgH.UName())
  }

  e.meter = otel.Metric(rpc, met)

  return e
}

// Finish ...
func (e *env) Finish() {
  defer e.Release()
  defer e.L.Flush()
  defer func() {
    if e.span != nil {
      e.span.End()
    }
    if e.meter != nil {
      e.meter.End()
    }
  }()

  if r := recover(); r != nil {
    // e.L.error("stack", string(debug.Stack()))
    if e.err == nil {
      e.err = fmt.Errorf("fail: unexpected panic")
    }

    fal := utils.GetPanicFrame(panicFrameSkip)
    e.L.CallerSkip(-1)
    // defer e.L.CallerSkip(envLogSkip)
    e.L.Fatal("S-F", path.Base(fal.Function),
      "S-pos", fmt.Sprintf("%s:%d", path.Base(fal.File), fal.Line),
      "panic", fmt.Sprintf("%+v", r))
  } else {
    e.L.Event("*", "finished")
  }

  for _, v := range e.Cf_ {
    v()
  }

  if e.enc == nil || e.dec == nil {
    return
  }

  if e.err != nil {
    e.msgH.Code = 1
    e.msgH.Reason = e.err.Error()
    e.ReV = nil
  }

  // e.L.error(e.err)
  e.L.CallerSkip(envLogSkip)
  e.writeRetV()
}

func (e *env) Return() {
  r := recover()
  if r == nil {
    e.L.Event("*", "returned")
    return
  }

  if e.err == nil {
    e.err = fmt.Errorf("fail: unexpected panic")
  }

  fal := utils.GetPanicFrame(panicFrameSkip + 1)
  e.L.CallerSkip(-1)
  defer e.L.CallerSkip(envLogSkip)
  e.L.Fatal("S-F", path.Base(fal.Function),
    "S-pos", fmt.Sprintf("%s:%d", path.Base(fal.File), fal.Line),
    "panic", fmt.Sprintf("%+v", r))
}

func (e *env) Release() {
  log.Release(e.L)
  if e.enc != nil {
    e.enc.Release()
  }
  if e.dec != nil {
    e.dec.Release()
  }
  pool.Put(e)
  return
}

func (e *env) Assert() {
  if e.err == nil {
    return
  }

  panic("")
}

func (e *env) HasError() bool {
  return e.err != nil
}

func (e *env) ResetErr() {
  e.err = nil
}

func (e *env) AssertErr(err error, clear ...func()) {
  if err == nil {
    return
  }

  for _, value := range clear {
    value()
  }
  e.L.Error("msg", err)

  e.err = err
  panic("")
}

func (e *env) AssertBool(ok bool, args ...any) {
  if ok {
    return
  }

  err := fmt.Errorf("fail:assertbool, info:%+v", args)
  e.L.Error("msg", err)
  e.err = err
  panic("")
}

func (e *env) Event(args ...interface{}) {
  e.L.Event(args...)
}

func (e *env) PrintParams(v ...reflect.Value) {
  if len(v) < 3 {
    return
  }

  str := []string{}
  for _, value := range v[2:] {
    s := fmt.Sprintf("%+v", value)
    str = append(str, string(s))
  }
  e.L.Event("params", str)
}

func (e *env) GetDec() decoder {
  return e.dec
}

func (e *env) GetMsgHeader() *Tjatse {
  if e.span != nil {
    e.msgH.SetSID(e.span.SID())
    e.msgH.SetTID(e.span.TID())
  }

  return e.msgH
}

func (e *env) SetReV(v []reflect.Value) {
  e.ReV = v
}

func (e *env) writeRetV() {
  if e.CheckErr() {
    return
  }

  for i := 0; i < len(e.ReV); i++ {
    err := e.enc.Encode(e.ReV[i].Interface())
    if err != nil {
      e.L.Event("msg", err)
      break
    }
  }
}

func (e *env) CheckErr(rtv ...reflect.Value) bool {
  // err := e.err
  var eErr *Err

  rtv_ := e.ReV
  if len(rtv) > 0 {
    rtv_ = rtv
  }

  for {
    if e.err != nil {
      break
    }

    if len(rtv_) == 0 {
      return false
    }

    ler := rtv_[len(rtv_)-1].Interface()
    if ler == nil {
      if len(rtv) == 0 {
        e.ReV = e.ReV[:len(e.ReV)-1]
      }
      return false
    }

    if err_, ok := ler.(*Err); ok {
      eErr = err_
      e.err = err_
    } else if err_, ok := ler.(error); ok {
      e.err = err_
    } else {
      return false
    }

    break
  }

  if len(rtv) > 0 {
    return true
  }

  if eErr != nil {
    e.msgH.Code = int32(eErr.Code)
  } else {
    e.msgH.Code = int32(InternalServerErr)
  }
  e.msgH.Reason = e.err.Error()
  e.enc.Encode(e.msgH)

  return true
}

type encoder interface {
  Encode(v any) error
  Release()
}
type decoder interface {
  Header() (*Tjatse, error)
  Decode(v any) error
  Release()
}

var (
  pool = sync.Pool{
    New: func() any {
      return &env{}
    },
  }
  envLogSkip     = 1
  panicFrameSkip = 3
)
