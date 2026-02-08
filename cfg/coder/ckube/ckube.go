package ckube

import (
  "fmt"
  "reflect"
  "strings"

  "github.com/yoophox/emei/cfg/source/inter"
  "github.com/yoophox/emei/cfg/values"
  "github.com/yoophox/emei/errs"
)

type values_ struct {
  s reflect.Value
}

// ...
func Encode(s inter.Source) (values.Values, error) {
  s_, _ := s.Read()
  v_ := reflect.ValueOf(s_)
  v_ = getElemOfPointer(v_)
  if v_.Kind() != reflect.Struct {
    return nil, errs.ErrorF("err.cfg.encoder.struct.encode", "source is not a structure value")
  }
  return &values_{v_}, nil
}

func (v *values_) Read(p string) (string, error) {
  if p == "" {
    return "", fmt.Errorf("fail: cfg encoder, encoder:kube, err:path is empty")
  }

  paths := strings.Split(p, ".")

  var ret reflect.Value = v.s
  for _, ps_ := range paths {
    if ret.Kind() != reflect.Struct {
      return "", errs.ErrorF("err.cfg.encoder.struct.read", "v.s is not struct")
    }
    ret := ret.FieldByName(ps_)
    if ret.Interface() == nil {
      return "", errs.ErrorF("err.cfg.encoder.struct.read.no", "not exist value for path:%s", p)
    }
  }

  return fmt.Sprintf("%v", ret.Interface()), nil
}

func (e *values_) Scan(p string, v any) (err error) {
  defer func() {
    if r := recover(); r != nil {
      err = fmt.Errorf("fail:cfg scan, coder:struct, reason:panic, msg:%+v, path:%s", r, p)
    }
  }()

  if p == "" {
    return fmt.Errorf("fail: cfg encoder, encoder:struct, err:path is empty")
  }

  v_ := reflect.ValueOf(v)
  if v_.Kind() != reflect.Pointer {
    return errs.ErrorF("err.cfg.encoder.struct.scan", "v is not pointer")
  }

  paths := strings.Split(p, ".")
  var ret reflect.Value = e.s
  for _, ps_ := range paths {
    ret = getElemOfPointer(ret)
    if ret.Kind() != reflect.Struct {
      return errs.ErrorF("err.cfg.encoder.struct.scan", "not exist value for path:%s, by:non-structured values", p)
    }
    ret := ret.FieldByName(ps_)
    if !ret.IsValid() {
      return errs.ErrorF("err.cfg.encoder.struct.scan", "not exist value for path:%s, by: no feild", p)
    }
  }

  return assign(ret, v_)
}

func (v *values_) Set(p string, o any) error {
  return nil
}

// getElemOfPointer ...
func getElemOfPointer(v reflect.Value) reflect.Value {
  if v.Kind() == reflect.Pointer {
    v = v.Elem()
  }

  return v
}
