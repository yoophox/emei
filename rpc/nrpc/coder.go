package nrpc

import (
  "bufio"
  "encoding/gob"
  "encoding/json"
  "fmt"
  "io"
  "net/rpc"
  "reflect"

  "github.com/yolksys/emei/env"
  "github.com/quic-go/quic-go"
)

type encoder struct {
  s *server
  i int
}

func (e *encoder) Encode(v any) error {
  if e.i == 0 {
    e.i++
    m, ok := v.(*env.Tjatse)
    if ok {
      e.s.res.MsgH = m
      return nil
    }
  }

  vs_, _ := json.Marshal(v)
  e.s.res.Data = append(e.s.res.Data, string(vs_))
  return nil
}

func (e *encoder) Release() {
}

// newEnc ...
func newEnc(s *server) *encoder {
  return &encoder{s: s}
}

type decoder struct {
  s *server
  i int
}

func (d *decoder) Decode(v any) error {
  vv_ := reflect.ValueOf(v)
  defer func() { d.i++ }()

  if vv_.Kind() != reflect.Pointer || vv_.IsNil() {
    return fmt.Errorf("fail:nrpc decode, err:kind, msg:v is not pointer or is nil, kind:%+v", vv_.Kind())
  }

  return json.Unmarshal([]byte(d.s.req.Data[d.i]), v)
}

func (d *decoder) Header() (*env.Tjatse, error) {
  err := d.s.gobServerCodec.ReadRequestBody(&d.s.req)
  if err != nil {
    return nil, err
  }

  return d.s.req.MsgH, nil
}

func (d decoder) Release() {
}

// newDec ...
func newDec(s *server) *decoder {
  return &decoder{s: s}
}

type gobServerCodec struct {
  rwc    io.ReadWriteCloser
  dec    *gob.Decoder
  enc    *gob.Encoder
  encBuf *bufio.Writer
  closed bool
}

func (c *gobServerCodec) ReadRequestHeader(r *rpc.Request) error {
  return c.dec.Decode(r)
}

func (c *gobServerCodec) ReadRequestBody(body any) error {
  return c.dec.Decode(body)
}

func (c *gobServerCodec) WriteResponse(r *rpc.Response, body any) (err error) {
  if err = c.enc.Encode(r); err != nil {
    if c.encBuf.Flush() == nil {
      // Gob couldn't encode the header. Should not happen, so if it does,
      // shut down the connection to signal that the connection is broken.
      // log.Println("rpc: gob error encoding response:", err)
      c.Close()
    }
    return
  }
  if err = c.enc.Encode(body); err != nil {
    if c.encBuf.Flush() == nil {
      // Was a gob problem encoding the body but the header has been written.
      // Shut down the connection to signal that the connection is broken.
      // log.Println("rpc: gob error encoding body:", err)
      c.Close()
    }
    return
  }
  return c.encBuf.Flush()
}

func (c *gobServerCodec) Close() error {
  if c.closed {
    // Only call c.rwc.Close once; otherwise the semantics are undefined.
    return nil
  }
  c.closed = true
  return c.rwc.Close()
}

// newSvrGobCodec ...
func newSvrGobCodec(conn io.ReadWriteCloser) *gobServerCodec {
  buf := bufio.NewWriter(conn)
  srv := &gobServerCodec{
    rwc:    conn,
    dec:    gob.NewDecoder(conn),
    enc:    gob.NewEncoder(buf),
    encBuf: buf,
  }

  return srv
}

type gobClientCodec struct {
  rwc    io.ReadWriteCloser
  dec    *gob.Decoder
  enc    *gob.Encoder
  encBuf *bufio.Writer
}

func (c *gobClientCodec) WriteRequest(r *rpc.Request, body any) (err error) {
  if err = c.enc.Encode(r); err != nil {
    return
  }
  if err = c.enc.Encode(body); err != nil {
    return
  }
  return c.encBuf.Flush()
}

func (c *gobClientCodec) ReadResponseHeader(r *rpc.Response) error {
  return c.dec.Decode(r)
}

func (c *gobClientCodec) ReadResponseBody(body any) error {
  return c.dec.Decode(body)
}

func (c *gobClientCodec) Close() error {
  return c.rwc.Close()
}

// newCliGobCodec ...
func newCliGobCodec(conn quic.Stream) rpc.ClientCodec {
  encBuf := bufio.NewWriter(conn)
  codec := &gobClientCodec{conn, gob.NewDecoder(conn), gob.NewEncoder(encBuf), encBuf}
  return codec
}
