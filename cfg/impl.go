package cfg

import (
  "time"

  "github.com/yolksys/emei/cfg/values"
)

type config struct {
  v values.Values
}

func (v *config) Bool(p string) (bool, error) {
  return false, nil
}

func (v *config) Int(p string) (int, error) {
  return 0, nil
}

func (v *config) String(p string) (string, error) {
  return "", nil
}

func (v *config) Float64(p string) (float64, error) {
  return 0.0, nil
}

func (v *config) Duration(p string) (time.Duration, error) {
  return time.Duration(0), nil
}

func (v *config) StringSlice(p string) ([]string, error) {
  return nil, nil
}

func (v *config) StringMap(p string) (map[string]string, error) {
  return map[string]string{}, nil
}

func (v *config) Scan(p string, ret any) error {
  return v.v.Scan(p, ret)
}

func (v *config) Bytes(p string) ([]byte, error) {
  return nil, nil
}
