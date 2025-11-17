package local

import (
  "errors"
  "fmt"
  "os"

  "github.com/yolksys/emei/cfg/source/inter"
)

// Load ...
func Load(path string) (inter.Source, error) {
  // f, err := os.OpenFile(path, os.O_WRONLY, 0777)
  f, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  return &local{
    f: f,
  }, nil
}

type local struct {
  f *os.File
}

func (s *local) Read() (any, error) {
  fs_, err := s.f.Stat()
  if err != nil {
    return nil, err
  }

  if fs_.Size() > int64(maxFileSzie) {
    return nil, errors.New(fmt.Sprintf("fail: local read,msg:size of file is big,maxszie: %v", maxFileSzie))
  }

  b := make([]byte, fs_.Size(), fs_.Size()+6)
  _, err = s.f.Read(b)
  if err != nil && err.Error() != "EOF" {
    return nil, err
  }

  return string(b), nil
}

func (s *local) Write(v any) error {
  return nil
}

var (
  maxFileSzie = 1024 * 1024
  p           = 0
)
