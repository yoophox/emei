package otel

import (
  "context"

  "go.opentelemetry.io/otel/attribute"
  // sdktrace "go.opentelemetry.io/otel/sdk/trace"
  "go.opentelemetry.io/otel/trace"
)

// global func
func newSpanCtxFromTraceData(d TracedData) context.Context {
  // sc_ := trace.SpanContextFromContext(context.Background())
  var scc trace.SpanContextConfig
  scc.TraceID, _ = trace.TraceIDFromHex(d.TraceId())
  scc.SpanID, _ = trace.SpanIDFromHex(d.TraceSpanId())
  scc.TraceFlags = 01
  scc.Remote = false
  return trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(scc))
}

// spanEnder ...
func spanEnder() {
  for {
    select {
    case s := <-_spanEndChan:
      s.End()
    }
  }
}

// --------------------------------------------------------
// for span
func (s *span) End() {
  //_spanEndChan <- s.s
  s.s.End()
}

func (s *span) SID() string {
  return s.s.SpanContext().SpanID().String()
}

func (s *span) TID() string {
  return s.s.SpanContext().TraceID().String()
}

func (s *span) AddAttri(key, value string) {
  s.s.SetAttributes(attribute.String(key, value))
}

// --------------------------------------------------------
// for meter
func (m *meter) End() {
  c, ok := _apiRealTimeCnt[m.key]
  if ok {
    c.Add(context.Background(), -1)
  }
}
