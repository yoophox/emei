package web

import (
  "bytes"
  "encoding/json"
  "io"
  "net/http"
  "sync"

  "github.com/yolksys/emei/env"
)

type writer struct {
  w io.Writer
}

func (i *writer) Write(buf []byte) (int, error) {
  return i.w.Write(buf)
}

type reader struct {
  r io.ReadCloser
}

func (i *reader) Read(buf []byte) (int, error) {
  return i.r.Read(buf)
}

type encoder struct {
  // r io.Reader
  w   *writer
  buf *bytes.Buffer
  enc *json.Encoder
  i   int
}

// newEnc ...
func newEnc(w ...io.Writer) *encoder {
  enc := poolEnc.Get().(*encoder)
  if len(w) == 1 {
    enc.w.w = w[0]
  } else {
    enc.w.w = enc.buf
  }
  enc.i = 0
  // enc.enc = json.NewEncoder(enc.w)
  return enc
}

func (e *encoder) Encode(v any) error {
  if e.i == 0 {
    e.i++
    m, ok := v.(*env.Tjatse)
    if ok {
      return e.enc.Encode(m)
    } else {
      err := e.enc.Encode(&struct{}{})
      if err != nil {
        return err
      }

    }
  }
  return e.enc.Encode(v)
}

func (e *encoder) getReader() io.Reader {
  return e.buf
}

func (e *encoder) Release() {
  e.w.w = nil
  e.buf.Reset()
  poolEnc.Put(e)
}

type decoder struct {
  r   *reader
  req *http.Request
  dec *json.Decoder
}

// newDec ...
func newDec(v any) *decoder {
  d := poolDec.Get().(*decoder)
  switch r := v.(type) {
  case *http.Request:
    d.req = r
    d.r.r = r.Body
  case io.ReadCloser:
    d.r.r = r
  }

  return d
}

func (d *decoder) Header() (*env.Tjatse, error) {
  h := &env.Tjatse{}
  ua_ := d.req.Header.Get("user-agent")
  if ua_ == defaultUserAgent {
    err := d.Decode(h)
    if err == io.EOF {
      return h, nil
    }
    return h, err
  }

  jwtc, err := d.req.Cookie("jwt")
  if err == nil {
    h.Jwt = jwtc.Value
  }
  return h, nil
}

func (d *decoder) Decode(v any) error {
  return d.dec.Decode(v)
}

func (d *decoder) Release() {
  d.r.r.Close()
  d.r.r = nil
  poolDec.Put(d)
}

var (
  poolBuf = sync.Pool{
    New: func() any {
      return bytes.NewBuffer([]byte{})
    },
  }

  poolwriter = sync.Pool{
    New: func() any {
      r := &writer{}
      return r
    },
  }

  poolEnc = sync.Pool{
    New: func() any {
      b := poolBuf.Get().(*bytes.Buffer)
      b.Reset()
      enc := &encoder{
        w:   poolwriter.Get().(*writer),
        buf: b,
      }
      enc.enc = json.NewEncoder(enc.w)

      return enc
    },
  }

  poolReader = sync.Pool{
    New: func() any {
      return &reader{}
    },
  }

  poolDec = sync.Pool{
    New: func() any {
      r := poolReader.Get().(*reader)
      d := &decoder{
        r:   r,
        dec: json.NewDecoder(r),
      }

      return d
    },
  }
)
