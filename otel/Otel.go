package otel

import (
  "context"

  "github.com/yolksys/emei/cfg"
  "github.com/yolksys/emei/utils"
  "go.opentelemetry.io/otel/log"
  "go.opentelemetry.io/otel/metric"
)

// Trace ...
func Trace(d TracedData) Span {
  if _tracer == nil {
    return nil
  }

  var ctx context.Context = context.Background()
  if d.TraceId() != "" {
    ctx = newSpanCtxFromTraceData(d)
  }
  _, s := _tracer.Start(ctx, cfg.Service)
  return &span{s}
}

// Metric ...
func Metric(key string, unit ...any) {
  if _meter == nil {
    return
  }

  var i int64 = 1
  var ok bool
  // i, ok = unit.(int64)
  if len(unit) > 0 {
    i, ok = unit[0].(int64)
    if !ok {
      return
    }
  }

  m, ok := _allMeters[key]
  if !ok {
    return
  }

  m.Add(context.TODO(), i)
}

// InitAllUpDownCounter ...
// called in init function of package
// @args: key, desc pair
func InitAllUpDownCounter(args ...string) {
  if _meter == nil {
    return
  }

  if !utils.IsCalledFromInit() {
    panic("this func must be called from init")
  }

  lastI := len(args) - 1
  for i := 0; i < len(args); i += 2 {
    if i == lastI {
      break
    }

    key := args[i]
    var err error
    _allMeters[key], err = _meter.Int64UpDownCounter(key, metric.WithDescription(args[i+1]))
    if err != nil {
      panic(err)
    }
    _allMetersCnt[key] = 0
  }
}

// Log ...
func Log(c *log.Record) {
  if _logger == nil {
    return
  }

  _logger.Emit(context.Background(), *c)
}
