package console

import (
  "context"
  "os"

  "github.com/yolksys/emei/log/backend/intra"
  "github.com/yolksys/emei/log/cache"
)

type console struct {
  c   chan *cache.LogRecord
  ctx context.Context
}

func New(ctx context.Context) intra.Backend {
  c := &console{
    c:   make(chan *cache.LogRecord, 1000),
    ctx: ctx,
  }

  go do(c.c)

  return c
}

func (c *console) Write(msg *cache.LogRecord) {
  c.c <- msg
}

func do(c chan *cache.LogRecord) {
  for {
    select {
    case m := <-c:
      os.Stdout.Write(m.Buf)
      cache.Put(m.Buf)
      m.Buf = nil
      cache.ReleaseLogRecord(m)
    }
  }
}

func init() {
  intra.RegisteredBck["console"] = New
}
