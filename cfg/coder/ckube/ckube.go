package ckube

import (
  "encoding/json"
  "fmt"
  "reflect"
  "strings"

  "k8s.io/api/core/v1"
  "github.com/yolksys/emei/cfg/source/inter"
  "github.com/yolksys/emei/cfg/values"
)

type values_ struct {
  s *v1.Service
}

// ...
func Encode(s inter.Source) (values.Values, error) {
  s_, _ := s.Read()
  return &values_{s_.(*v1.Service)}, nil
}

func (v *values_) Read(p string) (string, error) {
  if p == "" {
    return "", fmt.Errorf("fail: cfg encoder, encoder:kube, err:path is empty")
  }
  paths := strings.Split(p, ".")

  switch paths[0] {
  case "ports[*]":
    p, _ := json.Marshal(v.s.Spec.Ports)
    return string(p), nil
  case "metadata":
    if len(paths) == 3 && paths[1] == "labals" {
      l, ok := v.s.ObjectMeta.Labels[paths[2]]
      if ok {
        return l, nil
      }
    }
  }

  return "", fmt.Errorf("fail:cfg encoder, encoder:kube, err: path, path:%s", p)
}

func (e *values_) Scan(p string, v any) (err error) {
  defer func() {
    if r := recover(); r != nil {
      err = fmt.Errorf("fail:cfg scan, coder:kube, reason:panic, msg:%+v, path:%s", r, p)
    }
  }()

  if p == "" {
    return fmt.Errorf("fail: cfg encoder scan, encoder:kube, err:path is empty")
  }
  vt_ := reflect.TypeOf(v)
  if vt_.Kind() != reflect.Pointer {
    return fmt.Errorf("fail:cfg scan, coder:ckube, err:v is not pointer, v:%s", vt_)
  }

  paths := strings.Split(p, ".")

  switch paths[0] {
  case "ports[*]":
    // p, _ := json.Marshal(v.s.Spec.Ports)
    tt_ := reflect.MakeSlice(vt_.Elem(), len(e.s.Spec.Ports), len(e.s.Spec.Ports))
    for key, value := range e.s.Spec.Ports {
      tti := tt_.Index(key)
      ivn_ := tti.FieldByName("Name")
      ivp_ := tti.FieldByName("Port")
      if ivn_.IsNil() || ivp_.IsNil() {
        continue
      }
      ivn_.SetString(value.Name)
      ivp_.SetString(fmt.Sprintf("%d", value.Port))
    }
    vv_ := reflect.ValueOf(v)
    vv_.Elem().Set(tt_)
    return nil

  case "metadata":
    if len(paths) == 3 && paths[1] == "labals" {
      l, ok := e.s.ObjectMeta.Labels[paths[2]]
      if ok {
        return fmt.Errorf("fail:cfg scan, coder:kube, err:path, path:%s", p)
      }
      reflect.ValueOf(v).Elem().SetString(l)
      return nil
    }
  }

  return fmt.Errorf("fail:cfg skan, encoder:kube, err:path, path:%s", p)
}

func (v *values_) Set(p string, o any) error {
  return nil
}
