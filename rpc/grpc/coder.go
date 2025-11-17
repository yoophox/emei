package grpc

import (
  "encoding/json"
  "fmt"
  "reflect"

  "github.com/yolksys/emei/env"
)

type encoder struct {
  res *Response
  i   int
}

func (e *encoder) Encode(v any) error {
  if e.i == 0 {
    e.i++
    if m, ok := v.(*env.Tjatse); ok {
      e.res.MsgH = &Header{
        Code:   &m.Code,
        Reason: &m.Reason,
      }
      return nil
    }
  }
  vs_, _ := json.Marshal(v)
  e.res.ResValues = append(e.res.ResValues, string(vs_))
  return nil
}

func (e *encoder) Release() {
}

// newEnc ...
func newEnc(r *Response) *encoder {
  return &encoder{res: r}
}

type decoder struct {
  req *Request
  i   int
}

func (d *decoder) Decode(v any) error {
  vv_ := reflect.ValueOf(v)
  defer func() { d.i++ }()

  if vv_.Kind() != reflect.Pointer || vv_.IsNil() {
    return fmt.Errorf("fail:grpc decode, err:kind, msg:v is not pointer or is nil, kind:%+v", vv_.Kind())
  }

  return json.Unmarshal([]byte(d.req.ReqParams[d.i]), v)
}

func (d *decoder) Header() (*env.Tjatse, error) {
  h := &env.Tjatse{}
  _h := d.req.MsgH
  h.Mid = _h.Mid
  h.Jwt = _h.Jwt
  h.Sid = _h.Sid
  return h, nil
}

func (d decoder) Release() {
}

// newDec ...
func newDec(r *Request) *decoder {
  return &decoder{req: r}
}
