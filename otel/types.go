package otel

import (
  "go.opentelemetry.io/otel/attribute"
  "go.opentelemetry.io/otel/log"
  "go.opentelemetry.io/otel/metric"
  "go.opentelemetry.io/otel/trace"
)

type (
  LogRecord = log.Record
  LogKV     = log.KeyValue
  TraceId   = trace.TraceID
  AttriKV   = attribute.KeyValue
)

type TracedData interface {
  TraceId() string
  TraceSpanId() string
  // TraceFlag() string
  // TraceState() string
  SetSID(i string) // set span id. i must be hex string
  SetTID(i string) // set TraceId. i must be hex string
}

type Span interface {
  End()        // close()
  SID() string // span id
  TID() string // tracce id
  AddAttri(key, value string)
  // IsEnded() bool
}

type span struct {
  s trace.Span
}

type Meter interface {
  End()
}

type meter struct {
  rpc, met, key string
}

type resetableConter struct {
  cnt       uint64
  upDownCnt metric.Int64UpDownCounter
}
