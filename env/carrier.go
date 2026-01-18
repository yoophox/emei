package env

import "io"

type defaultCarrier struct {
  io.ReadWriter
}

func (d *defaultCarrier) Inject(e *Tjatse) error {
  return nil
}

func (d *defaultCarrier) extract(e *Tjatse) error {
  return nil
}

// ...
func newDftCarrier(io_ io.ReadWriter) *defaultCarrier {
  return nil
}
