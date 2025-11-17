package yaml

import (
  "fmt"
  "io"
  "strings"
  "sync"

  yaml "github.com/goccy/go-yaml"
  "github.com/yolksys/emei/cfg/source/inter"
  "github.com/yolksys/emei/cfg/values"
)

// Encoder ...
func Encode(s inter.Source) (values.Values, error) {
  str_, e := s.Read()
  if e != nil {
    return nil, e
  }

  str := str_.(string)
  c := &yamlValues{
    raw:     str,
    pReader: map[string]*yaml.Path{},
    r:       strings.NewReader(str),
  }
  c.init()

  return c, nil
}

type yamlValues struct {
  raw     string
  r       io.Reader
  pReader map[string]*yaml.Path
  m       sync.RWMutex
}

func (y *yamlValues) init() error {
  // y.r = strings.NewReader(y.raw)
  // y.r.
  return nil
}

func (y *yamlValues) Read(p string) (string, error) {
  p = "$." + p
  y.m.RLock()
  pr_, ok := y.pReader[p]
  y.m.RUnlock()
  if !ok {
    var e error
    pr_, e = yaml.PathString(p)
    if e != nil {
      return "", fmt.Errorf("fail:cfg read, coder: yaml, reason:PathString, err: %s, path:%s", e.Error(), p)
    }
    y.m.Lock()
    y.pReader[p] = pr_
    y.m.Unlock()
  }

  nd_, e := pr_.ReadNode(strings.NewReader(y.raw))
  if e != nil {
    return "", fmt.Errorf("fail:cfg read, coder: yaml, reason:ReadNode, err: %s, path:%s", e.Error(), p)
  }

  return nd_.String(), nil
}

func (y *yamlValues) Scan(p string, v any) error {
  s, err := y.Read(p)
  if err != nil {
    return err
  }

  return yaml.Unmarshal([]byte(s), v)
}

func (y *yamlValues) Set(p string, v any) error {
  return nil
}
