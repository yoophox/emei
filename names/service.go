package names

const (
  // @@ = AT_, _AT if can't use @
  NAME_SERVICE_SELF           = "@@self"
  NAME_SERVICE_OTEL_COLLECTOR = "@@otel-collector"
  NAME_SERVICE_PKI            = "@@pki"
  NAME_SERVICE_UUID           = "@@uuid@" // don't use "uuid" as default service name
  NAME_SERVICE_VALKEY         = "@@valkey"
)

const (
  // annotations: @@service-quic
  NAME_SERVICE_PORT_DEFAULT_QUIC  = "quic"
  NAME_SERVICE_PORT_DEFAULT_TCP   = "tcp"
  NAME_SERVICE_PORT_OTEL_COLLETOR = "grpc"
  NAME_SERVICE_PORT_VALKEY_CLIENT = "clt"
)
