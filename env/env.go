package env

import (
  "context"
  "fmt"
  "path"
  "sync"

  "github.com/yoophox/emei/errs"
  "github.com/yoophox/emei/jwt"
  "github.com/yoophox/emei/log"
  "github.com/yoophox/emei/log/core"
  "github.com/yoophox/emei/otel"
  "github.com/yoophox/emei/utils"
)

type env struct {
  log.Logger
  // Cf_ []func()        // callback
  // Ecf []func()        // error call back
  // ReV []reflect.Value // return values
  par  *env    // up env
  tjax *Tjatse // trace info for oepntele

  span otel.Span
  // rpc   string
  // met   string
  // meter otel.Meter
  err    error
  jwt    jwt.JWT
  wgx    sync.WaitGroup
  ctx    context.Context
  isDone bool
}

// New ...
func new(tja *Tjatse) *env {
  e := pool.Get().(*env)
  e.Logger = log.New(context.Background(), core.WithCacheMode())
  e.Logger.CallerSkip(envLogSkip)
  e.Logger.Event("*", "start")
  e.wgx.Add(1)

  //
  if tja != nil {
    e.tjax = tja
  } else {
    e.tjax = &Tjatse{}
  }
  e.span = otel.Trace(e.tjax)
  if e.span == nil {
    e.AddAttri("uid", e.uid())
    e.AddAttri("uname", e.uname())
  } else {
    e.Logger.SetTraceId(e.span.TID())
    e.Logger.AddAttri("spanid", e.span.SID())
  }

  return e
}

// Finish ...
func (e *env) Finish() {
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
}

func (e *env) Trace() {
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
  e.Logger.Flush()
  log.Release(e.Logger)
  pool.Put(e)
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

func (e *env) AssertErr(err error, errid ...errs.ErrId) {
  if err == nil {
    return
  }

  if len(errid) == 1 {
    err = errs.Wrap(err, errid[0])
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

  err := fmt.Errorf(fmt_, args...)
  e.Logger.Error("msg", err)
  e.err = errs.Wrap(err, eid)
  panic("")
}

func (e *env) Clone() *env {
  var e_ env
  var tja Tjatse
  tja = *e.tjax
  e.tjax = &tja
  e.par = e
  e.span = otel.Trace(&tja)
  e.Logger = log.New(context.Background(), log.WithCacheMode)
  return &e_
}

func (e *env) Err() error {
  return e.err
}

func (e *env) TID() string {
  if e.span != nil {
    return e.span.TID()
  }

  return ""
}

func (e *env) JWT(j ...any) jwt.JWT {
  if len(j) == 1 {
    switch j_ := j[0].(type) {
    case string:
      e.jwt = jwt.Parse(j_)
      e.tjax.Jwt = j_
    case jwt.JWT:
      e.jwt = j_
    }
  }

  if e.jwt == nil {
    e.jwt = jwt.Parse(e.tjax.Jwt)
  }

  return e.jwt
}

func (e *env) Done() <-chan struct{} {
  return e.ctx.Done()
}

func (e *env) IsDone() bool {
  return e.ctx.Err() != nil
}

func (e *env) Go(f any, args ...any) {
  e.wgx.Add(1)
  go func() {
    e_ := e.Clone()
    defer func() {
      e_.Release()
      e.wgx.Done()
    }()
    // f(args...)
  }()
}

func (e *env) Wait() {
  e.wgx.Wait()
}

func (e *env) WaitAny() {
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
