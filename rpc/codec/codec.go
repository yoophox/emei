package codec

type CodecIx interface {
  Decode(v any) error
  Encode(v any) error
  // DecodeTyps(typs ...reflect.Type) ([]reflect.Value, error)
  // EncodeValues(vs_ ...reflect.Value) error
}

const (
  CODEC_GOB byte = iota
  CODEC_HTTP_JSON
  CODEC_GRPC
)
