package otel

import (
  "context"
  "fmt"
  "strings"

  "go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
  "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
  "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
  "go.opentelemetry.io/otel/sdk/log"
  "go.opentelemetry.io/otel/sdk/metric"
  "go.opentelemetry.io/otel/sdk/trace"

  "github.com/yolksys/emei/cfg"
  "github.com/yolksys/emei/kube"
  "github.com/yolksys/emei/utils"
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

  _totalCnt, err = _meter.Int64Counter("total_call_counter")
  if err != nil {
    panic(err)
  }
  _nrpcCnt, err = _meter.Int64Counter("nrpc_call_counter")
  if err != nil {
    panic(err)
  }
  _webCnt, err = _meter.Int64Counter("web_call_counter")
  if err != nil {
    panic(err)
  }
  _grpcCnt, err = _meter.Int64Counter("grpc_call_counter")
  if err != nil {
    panic(err)
  }

  // go spanEnder()
}

// createExports ...
func createExports() {
  var cor string
  err := cfg.GetCfgItem("otel.collector", &cor)
  if len(cor) > 9 && cor[:4] == "grpc" {
    g := strings.Split(cor, "://")
    if len(g) != 2 {
      panic(cor)
    }

    _, ip, port, err := kube.Lookup(g[1], "grpc")
    if err != nil {
      panic(err)
    }
    if !utils.IsIpv4(ip) {
      ip = "[" + ip + "]"
    }

    url := "http://" + ip + ":" + port
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
    return
  }

  fmt.Println("otel init error", err, cor)
  // panic("just support colletor of grpc: " + cor)
}
