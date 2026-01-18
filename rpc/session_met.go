package rpc

import (
  "fmt"
  "reflect"

  "github.com/yolksys/emei/env"
)

func (c *session) decode(typs ...reflect.Type) ([]reflect.Value, error) {
  ret := []reflect.Value{}
  for i, typ := range typs {
    var v reflect.Value

    if typ.Kind() == reflect.Pointer {
      v = reflect.New(typ.Elem())
    } else {
      v = reflect.New(typ)
    }

    err := c.Decode(v)
    if err != nil {
      c.zero(&ret, typs[i:]...)
      return ret, fmt.Errorf("type: %s,info: %s", typ.Name(), err.Error())
    }

    if typ.Kind() != reflect.Pointer {
      v = v.Elem()
    }
    ret = append(ret, v)
  }
  return ret, nil
}

func (c *session) encode(values ...any) error {
  for _, v := range values {
    err := c.Encode(v)
    if err != nil {
      return fmt.Errorf("value: %v, %v", v, err)
    }
  }

  return nil
}

func (c *session) zero(ret *([]reflect.Value), typs ...reflect.Type) {
  for _, typ := range typs {
    *ret = append(*ret, reflect.Zero(typ))
  }
}

func (c *session) close() {
  c.ReadWriteCloser.Close()
}

func (c *session) Inject(j *env.Tjatse) error {
  return nil
}

func (c *session) Extract(j *env.Tjatse) error {
  return nil
}
