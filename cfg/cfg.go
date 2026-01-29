package cfg

import (
  "errors"
  "fmt"
  "path"
  "strings"
  "time"

  "github.com/yoophox/emei/cfg/coder"
  "github.com/yoophox/emei/cfg/source"
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

func BuildCfgURI(source, path string, coder ...string) string {
  if coder != nil {
    return fmt.Sprintf("%s~%s~%s", source, path, coder[0])
  } else {
    return fmt.Sprintf("%s~%s", source, path)
  }
}

// getCoderType ...
func getCoderType(p string) string {
  ext := path.Ext(p)
  switch ext {
  case ".json":
    return CFG_CODER_JSON
  case ".yaml":
    return CFG_CODER_YAML
  default:
    return ""
  }
}
