package call

type Decoder interface {
  Decode(v any) error
}
