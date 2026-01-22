package codec

import (
  "net/http"
  "reflect"
  // "golang.org/x/net/http2"
)

// newHttp1_2 ...
//func NewHttp1_2(io_ io.ReadWriteCloser) error {
//  const preface = "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"
//  b := make([]byte, len(preface))
//  if _, err := io.ReadFull(io_, b); err != nil {
//    return err
//  }
//  if string(b) != preface {
//  }
//
//  framer := http2.NewFramer(io_, io_)
//  frame, err := framer.ReadFrame()
//  fmt.Println(frame, err)
//  return err
//}

func NewHttpJsonCodec(w http.ResponseWriter, r *http.Request)

type httpJsonCodec struct {
  w http.ResponseWriter
  r *http.Request
}

func (h *httpJsonCodec) Decode(v any) error {
  return nil
}

func (h *httpJsonCodec) Encode(v any) error {
  return nil
}

func (h *httpJsonCodec) DecodeTyps(typs ...reflect.Type) error {
  return nil
}

func (h *httpJsonCodec) Encodealues(vs_ ...reflect.Value) error {
  return nil
}
