package nrpc

import (
  "errors"
  "fmt"
  "io"
  "os"
  "time"

  "github.com/yolksys/emei/env"
  "github.com/quic-go/quic-go"
)

type stream struct {
  s       quic.Stream
  close   func()
  err     error
  leftLen uint16
}

func (c *stream) Write(b []byte) (int, error) {
  // do_ ...
  tol_ := 0

  do_ := func(b []byte) (int, error) {
    l := len(b)
    if l == 0 {
      return l, nil
    }
    err := c.writeLen(uint16(l))
    if err != nil {
      return 0, err
    }

    ll_ := 0
    for ll_ < l {
      p := b[ll_:]
      n, err := c.s.Write(p)
      ll_ += n
      if err != nil {
        return ll_, err
      }
      // b = b[n:]
    }

    tol_ += ll_

    return l, nil
  }

  c.s.SetWriteDeadline(time.Now().Add(30 * time.Second)) //
  l := len(b)
  for len(b) > _streamMaxWriteLen {
    n, err := do_(b[:_streamMaxWriteLen])
    if err != nil {
      c.err = err
      return n + tol_, err
    }
    b = b[_streamMaxWriteLen:]
  }

  n, err := do_(b)
  if err != nil {
    c.err = err
    return n + tol_, err
  }

  return l, nil
}

// Read ...
func (c *stream) Read(b []byte) (int, error) {
  l := len(b)
  ll_ := 9
  c.s.SetReadDeadline(time.Now().Add(30 * time.Second))
  if c.leftLen == 0 {
    var err error
    c.leftLen, err = c.readLen()
    if err != nil {
      if errors.Is(err, os.ErrDeadlineExceeded) {
        return ll_, nil
      }
      if errors.Is(err, _errSvrResed) {
        return ll_, err
      }
      c.err = err
      return ll_, err
    }
  }

  rLen := ll_ + int(c.leftLen)
  if l-ll_ < int(c.leftLen) {
    rLen = l
  }

  n, err := c.s.Read(b[ll_:rLen])
  ll_ += n
  if err != nil {
    if errors.Is(err, os.ErrDeadlineExceeded) {
      return ll_, nil
    }
    c.err = err
    return ll_, err
  }
  c.leftLen -= uint16(n)

  return ll_, nil
}

func (c *stream) Close() error {
  if c.err != nil && c.err != io.EOF {
    return nil
  }
  c.s.SetDeadline(time.Now().Add(30 * time.Second))
  l := 0
  for l < 2 {
    n, err := c.s.Write(_streamEofData)
    if err != nil {
      return err
    }

    l += n
  }

  if c.err == nil {
    buf := make([]byte, 2)
    for {
      _, err := c.Read(buf)
      if err != nil {
        if err == io.EOF {
          break
        }

        return err
      }
    }
  }

  c.err = nil
  c.leftLen = 0
  if c.close != nil {
    c.close()
  }
  return nil
}

func (c *stream) writeLen(l uint16) error {
  blen := 2
  b := []byte{byte(l >> 8), byte(l & 255)}
  for blen > 0 {
    b = b[len(b)-blen:]
    n, err := c.s.Write(b)
    if err != nil {
      if err == io.EOF {
        return fmt.Errorf("fail:nrpc stream writeLen, reason:eof")
      }
      return err
    }

    blen = blen - n
  }

  return nil
}

func (c *stream) readLen() (l uint16, err error) {
  r := func(l int) (b []byte, err error) {
    b = make([]byte, l)
    c.s.SetReadDeadline(time.Now().Add(30 * time.Second))
    l__ := 0
    for l__ < l {
      n, err := c.s.Read(b[l__:])
      if err != nil {
        if err == io.EOF {
          return nil, fmt.Errorf("fail:nrpc stream readLen, reason:eof")
        }
        return nil, err
      }
      l__ += n
    }
    return
  }

  b, err := r(2)
  if err != nil {
    return 0, err
  }

  if _streamEofData[0] == b[0] && _streamEofData[1] == b[1] {
    return 0, io.EOF
  }

  if b[0] == 255 {
    b, err = r(int(b[1]))
    if err != nil {
      return 0, err
    }

    return 0, fmt.Errorf("fail:svrhres, reason:%s, par:{%w}", string(b), _errSvrResed)
  }

  l = uint16(b[0])<<8 + uint16(b[1])

  return
}

func (c *stream) writeErr(e env.Env, err error) error {
  estr := err.Error()
  if len(estr) > 254 {
    estr = estr[:254]
  }
  d := []byte{255, byte(len(estr))}
  c.s.SetWriteDeadline(time.Now().Add(30 * time.Second))
  n, err := c.s.Write(d)

  if n != 2 || err != nil {
    e.Event("nrpc writeerr header", err.Error())
    return err
  }

  n, err = c.Write([]byte(estr))
  if n != len(estr) || err != nil {
    e.Event("nrpc writeerr estr", err.Error())
    return err
  }

  return nil
}
