package otel

import (
  "context"

  "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
  "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
  "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
  "go.opentelemetry.io/otel/sdk/log"
  "go.opentelemetry.io/otel/sdk/metric"
  "go.opentelemetry.io/otel/sdk/trace"

  "github.com/yoophox/emei/kube"
  "github.com/yoophox/emei/utils"
  "go.opentelemetry.io/otel/sdk/resource"
)

func init() {
  var err error

  _resource, err = resource.New(context.Background(), resource.WithContainerID(), resource.WithHostID())
  if err != nil {
    panic(err)
  }

  createExports()
  if _otelTraceExpoter == nil {
    return
  }

  _logProvider = log.NewLoggerProvider(log.WithProcessor(log.NewBatchProcessor(_otelLogExporter)))
  _metricProvider = metric.NewMeterProvider(metric.WithReader(metric.NewPeriodicReader(_otelMetricExporter)))
  _traceProvider = trace.NewTracerProvider(trace.WithBatcher(_otelTraceExpoter))

  _logger = _logProvider.Logger(_scope + "logger")
  _meter = _metricProvider.Meter(_scope + "metric")
  _tracer = _traceProvider.Tracer(_scope + "trace")
  // go spanEnder()
}

// createExports ...
func createExports() {
  s, err := kube.LookupServer("@@opentelemetry")
  if err != nil {
    panic(err)
  }

  ip := s.IP
  if !utils.IsIpv4(ip) {
    ip = "[" + ip + "]"
  }

  url := "http://" + ip + ":" + s.Port
  _otelLogExporter, err = otlploggrpc.New(context.Background(),
    otlploggrpc.WithEndpointURL(url))
  if err != nil {
    panic(err)
  }
  _otelMetricExporter, err = otlpmetricgrpc.New(context.Background(),
    otlpmetricgrpc.WithEndpointURL(url))
  if err != nil {
    panic(err)
  }
  _otelTraceExpoter, err = otlptracegrpc.New(context.Background(),
    otlptracegrpc.WithEndpointURL(url))
  if err != nil {
    panic(err)
  }

  // panic("just support colletor of grpc: " + cor)
}
