package env

import (
  "context"
  "fmt"
  "path"
  "reflect"
  "sync"

  "github.com/yolksys/emei/errs"
  "github.com/yolksys/emei/jwt"
  "github.com/yolksys/emei/log"
  "github.com/yolksys/emei/log/core"
  "github.com/yolksys/emei/otel"
  "github.com/yolksys/emei/utils"
)

type env struct {
  log.Logger
  // Cf_ []func()        // callback
  // Ecf []func()        // error call back
  // ReV []reflect.Value // return values
  Par  *env   // up env
  tjax Tjatse // trace info for oepntele

  span otel.Span
  // rpc   string
  // met   string
  // meter otel.Meter
  err     error
  actions []Action
  jwt     jwt.JWT
}

// New ...
func new(crr ...Carrier) *env {
  e := pool.Get().(*env)
  e.Logger = log.New(context.Background(), core.WithCacheMode())
  e.Logger.CallerSkip(envLogSkip)
  e.Logger.Event("*", "start")

  if e.actions == nil {
    e.actions = []Action{}
  } else {
    e.actions = e.actions[:0]
  }

  //
  if crr != nil {
    err := crr[0].Extract(&e.tjax)
    if err != nil {
      e.err = errs.Wrap(err, ERR_ID_ENV_EXTRACT)
      return e
    }
  } else {
    e.tjax.Mid = ""
  }

  e.span = otel.Trace(&e.tjax)
  if e.span == nil {
    e.AddAttri("uid", e.uid())
    e.AddAttri("uname", e.uname())
    return e
  }
  e.Logger.SetTraceId(e.span.TID())

  return e
}

// Finish ...
func (e *env) Finish(actf ...Action) {
  if actf != nil {
    for _, f := range actf {
      e.actions = append(e.actions, f)
    }

    return
  }
  defer e.Release()
  defer e.Logger.Flush()
  defer func() {
    if e.span != nil {
      e.span.End()
    }
    e.Logger.Event("*", "finished")
  }()

  if r := recover(); r != nil {
    // e.L.error("stack", string(debug.Stack()))
    if e.err == nil {
      e.err = fmt.Errorf("fail: unexpected panic")
    }

    fal := utils.GetCallerFrame(panicFrameSkip)
    e.Logger.CallerSkip(-1)
    // defer e.L.CallerSkip(envLogSkip)
    e.Logger.Fatal("S-F", path.Base(fal.Function),
      "S-pos", fmt.Sprintf("%s:%d", path.Base(fal.File), fal.Line),
      "panic", fmt.Sprintf("%+v", r))

    e_, ok := e.err.(*errs.Err)
    if ok {
      e.tjax.Code = string(e_.Eid)
    } else {
      e.tjax.Code = "err.default"
    }
    e.tjax.Reason = e.err.Error()
  }

  e.Logger.CallerSkip(envLogSkip)

  if e.err != nil || len(e.actions) > 0 {
    err := e.Propagate()
    if err != nil {
      e.Error("env propagate", err.Error())
    }
  }

  for _, f := range e.actions {
    err := f()
    if err != nil {
      e.Error("env action:"+reflect.TypeOf(f).Name(), err.Error())
      break
    }
  }
}

func (e *env) Return() {
  r := recover()
  if r == nil {
    e.Logger.Event("*", "returned")
    return
  }

  if e.err == nil {
    e.err = fmt.Errorf("fail: unexpected panic")
  }

  fal := utils.GetCallerFrame(panicFrameSkip + 1)
  e.Logger.CallerSkip(-1)
  defer e.Logger.CallerSkip(envLogSkip)
  e.Logger.Fatal("S-F", path.Base(fal.Function),
    "S-pos", fmt.Sprintf("%s:%d", path.Base(fal.File), fal.Line),
    "panic", fmt.Sprintf("%+v", r))
}

func (e *env) Release() {
  log.Release(e.Logger)
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

func (e *env) AssertErr(err error) {
  if err == nil {
    return
  }

  // for _, value := range clear {
  //   switch v := value.(type) {
  //   case ClearFunc:
  //     v()
  //   case errs.ErrId:
  //     err = errs.Wrap(err, v)
  //   }
  // }
  e.Logger.Error("msg", err)

  e.err = err
  panic(err)
}

func (e *env) AssertBool(ok bool, eid errs.ErrId, fmt_ string, args ...any) {
  if ok {
    return
  }

  err := fmt.Errorf(fmt_, args)
  e.Logger.Error("msg", err)
  e.err = errs.Wrap(err, eid)
  panic("")
}

func (e *env) Propagate(crr ...Carrier) error {
  if e.span != nil {
    e.tjax.SetSID(e.span.SID())
    e.tjax.SetTID(e.span.TID())
  }

  return nil
}

func (e *env) uid() string {
  if e.jwt == nil {
    e.jwt = jwt.Parse(e.tjax.Jwt)
  }

  return e.jwt.GetClaim(jwt.COMMON_USER_CLAIM_UID)
}

func (e *env) uname() string {
  if e.jwt == nil {
    e.jwt = jwt.Parse(e.tjax.Jwt)
  }

  return e.jwt.GetClaim(jwt.COMMON_USER_CLAIM_UNAME)
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
