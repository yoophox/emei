package values

type Values interface {
  Read(path string) (string, error)
  Set(path string, v any) error
  Scan(path string, v any) error
}
