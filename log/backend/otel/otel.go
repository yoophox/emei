package otel

import (
  "context"
  "time"

  "github.com/yolksys/emei/log/backend/intra"
  "github.com/yolksys/emei/log/cache"
  rotel "github.com/yolksys/emei/otel"
)

type otel struct {
  c   chan *cache.LogRecord
  ctx context.Context
}

// New ...
func New(ctx context.Context) intra.Backend {
  b := &otel{
    ctx: ctx,
    c:   make(chan *cache.LogRecord, 1000),
  }
  go do(b.c)
  return b
}

func (b *otel) Write(msg *cache.LogRecord) {
  b.c <- msg
}

// do ...
func do(c chan *cache.LogRecord) {
  for {
    select {
    case m := <-c:
      lr_ := &rotel.LogRecord{}
      lr_.SetBody(rotel.LogStringValue(string(m.Buf)))
      lr_.SetTimestamp(time.Now())
      // t, _ := rotel.LogTraceIdFromHex(m.TraceId)
      // lr_.SetTraceID(t)
      lr_.AddAttributes(rotel.LogKV{
        Key:   "traceID",
        Value: rotel.LogStringValue(m.TraceId),
      })
      for key, value := range m.Attris {
        lr_.AddAttributes(rotel.LogKV{
          Key:   key,
          Value: rotel.LogStringValue(value),
        })
      }
      rotel.Log(lr_)
    }
  }
}

func init() {
  intra.RegisteredBck["otel"] = New
}
