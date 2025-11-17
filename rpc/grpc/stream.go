package grpc

type grpcstm interface {
  Send(*CData) error
  Recv() (*CData, error)
}

type stream struct {
  c grpcstm
}

func (s *stream) Write(b []byte) (n int, err error) {
  return 0, nil
}

func (s *stream) Read(b []byte) (n int, err error) {
  return 0, nil
}

func (s *stream) Close() error {
  if c, ok := s.c.(Grpc_RWStreamClient); ok {
    return c.CloseSend()
  }
  return nil
}
