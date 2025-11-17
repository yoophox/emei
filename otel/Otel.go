package otel

import (
  "context"
  "fmt"

  "go.opentelemetry.io/otel/log"
)

// Trace ...
func Trace(d TracedData, name string) Span {
  if _tracer == nil {
    return nil
  }

  var ctx context.Context = context.Background()
  if d.TraceId() != "" {
    ctx = newSpanCtxFromTraceData(d)
  }
  _, s := _tracer.Start(ctx, name)
  return &span{s}
}

// Metric ...
func Metric(rpc, api string) Meter {
  if _meter == nil {
    return nil
  }

  ctx := context.Background()
  key := rpc + "." + api
  c, ok := _apiCnt[key]
  if ok {
    c.Add(ctx, 1)
  }
  c0, ok := _apiRealTimeCnt[key]
  if ok {
    c0.Add(ctx, 1)
  }
  _totalCnt.Add(ctx, 1)

  switch rpc {
  case "web":
    _webCnt.Add(ctx, 1)
  case "nrpc":
    _nrpcCnt.Add(ctx, 1)
  case "grpc":
    _grpcCnt.Add(ctx, 1)
  }

  return &meter{rpc, api, key}
}

// Log ...
func Log(c *log.Record) {
  if _logger == nil {
    return
  }

  _logger.Emit(context.Background(), *c)
}

// ...
func AddApiMeter(rpc, api string) {
  key := rpc + "." + api
  var err error
  _apiCnt[key], err = _meter.Int64Counter(key)
  if err != nil {
    fmt.Println("addapimeter err: ", err)
  }

  _apiRealTimeCnt[key], _ = _meter.Int64UpDownCounter(key + "_real-time")
}
