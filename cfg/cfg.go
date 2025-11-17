package cfg

import (
  "errors"
  "fmt"
  "path"
  "strings"
  "time"

  "github.com/yolksys/emei/cfg/coder"
  "github.com/yolksys/emei/cfg/source"
)

// p = "m.k.l"
type Config interface {
  Bool(p string) (bool, error)
  Int(p string) (int, error)
  String(p string) (string, error)
  Float64(p string) (float64, error)
  Duration(p string) (time.Duration, error)
  StringSlice(p string) ([]string, error)
  StringMap(p string) (map[string]string, error)
  Scan(p string, v any) error
  Bytes(p string) ([]byte, error)
}

// New ...
// uri: "source~url[~coder]"
func New(uri string) (Config, error) {
  sec := strings.Split(uri, "~")
  l := len(sec)
  if l != 2 && l != 3 {
    return nil, errors.New("Fail:cfg new,msg:uri formate:'source~url[~coder]',uri:" + uri)
  }

  cTyp := ""
  if l == 2 {
    cTyp = getCoderType(sec[1])
  } else {
    cTyp = sec[2]
  }

  s, err := source.Load(sec[0], sec[1])
  if err != nil {
    return nil, err
  }

  v, err := coder.Encode(cTyp, s)
  if err != nil {
    return nil, err
  }

  return &config{
    v: v,
  }, nil
}

// GetCfgItem ...
func GetCfgItem(p string, v any) error {
  if _cfg == nil {
    return fmt.Errorf("fail:getcfgitem, reason:have no default cfg file")
  }
  return _cfg.Scan(p, v)
}

// SetUID ...
func SetUID(uid string) {
  _uid = uid
}

// GetUID ...
func GetUID() string {
  return _uid
}

// getCoderType ...
func getCoderType(p string) string {
  ext := path.Ext(p)
  switch ext {
  case "cfg":
    return "cfg"
  case "yaml":
    return "yaml"
  default:
    return ""
  }
}
