package otel

import (
  "go.opentelemetry.io/otel/log"
  "go.opentelemetry.io/otel/metric"
  sdklog "go.opentelemetry.io/otel/sdk/log"
  sdkmetric "go.opentelemetry.io/otel/sdk/metric"

  "go.opentelemetry.io/otel/sdk/resource"
  sdktrace "go.opentelemetry.io/otel/sdk/trace"
  "go.opentelemetry.io/otel/trace"
)

var (
  _resource *resource.Resource

  _otelLogExporter    sdklog.Exporter
  _otelTraceExpoter   sdktrace.SpanExporter
  _otelMetricExporter sdkmetric.Exporter

  _logProvider    log.LoggerProvider
  _metricProvider metric.MeterProvider
  _traceProvider  trace.TracerProvider

  _logger log.Logger
  _meter  metric.Meter
  _tracer trace.Tracer
)

var (
  _totalCnt       metric.Int64Counter
  _nrpcCnt        metric.Int64Counter
  _webCnt         metric.Int64Counter
  _grpcCnt        metric.Int64Counter
  _apiCnt         = map[string]metric.Int64Counter{}
  _apiRealTimeCnt = map[string]metric.Int64UpDownCounter{}
)

var (
  _spanEndChan = make(chan trace.Span, 1)
  _scope       = "gighub.com/emei/"
)

var (
  LogStringValue    = log.StringValue
  LogTraceIdFromHex = trace.TraceIDFromHex
)
