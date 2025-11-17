package inter

type Reader interface {
  Read() (any, error)
}

type Writer interface {
  Write(n any) error
}

type Source interface {
  Reader
  Writer
}
