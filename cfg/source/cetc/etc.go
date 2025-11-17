package cetc

import (
  "fmt"
  "strings"

  "github.com/yolksys/emei/cfg/source/inter"
)

type etcs struct {
  s string
}

// Load ...
func Load(p string) (inter.Source, error) {
  p_ := strings.Split(p, "/")
  if len(p_) != 2 {
    return nil, fmt.Errorf("fail:cfg source load, source:etc, path:%s", p)
  }
  f, ok := RegedEtcFunc[p_[0]]
  if !ok {
    return nil, fmt.Errorf("fail:cfg source load, source:etc, err:type, path:%s", p)
  }
  s, err := f(p_[1])
  if err != nil {
    return nil, fmt.Errorf("fail: cfg source load, source: etc, par{%s}", err.Error())
  }
  return &etcs{s}, nil
}

func (s *etcs) Read() (any, error) {
  return s.s, nil
}

func (s *etcs) Write(v any) error {
  return nil
}

type etcLoad func(string) (string, error)

var RegedEtcFunc = map[string]etcLoad{}
