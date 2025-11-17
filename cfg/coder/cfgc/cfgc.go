package cfgc

import (
  "bytes"
  "encoding/json"
  "fmt"
  "unicode"

  "github.com/yolksys/emei/cfg/source/inter"
  "github.com/yolksys/emei/cfg/values"
  "github.com/tidwall/gjson"
)

// Encoder ...
func Encode(s inter.Source) (values.Values, error) {
  str_, e := s.Read()
  if e != nil {
    return nil, e
  }

  str := str_.(string)
  c := &cfgValues{
    raw: str,
    buf: bytes.NewBuffer([]byte(str)),
  }

  e = c.init()
  return c, e
}

type cfgValues struct {
  j   string
  raw string
  buf *bytes.Buffer
}

func (c *cfgValues) Read(p string) (string, error) {
  v := gjson.Get(c.j, p)
  s := v.Raw
  if s == "" && !v.Exists() {
    return "", fmt.Errorf("fail:cfgc read, reason:path not exist, path:%s", p)
  }
  return v.Raw, nil
}

func (c cfgValues) Scan(p string, v any) error {
  ss_, err := c.Read(p)
  if err != nil {
    return err
  }

  return json.Unmarshal([]byte(ss_), v)
}

func (c *cfgValues) Set(p string, v any) error {
  return nil
}

func (c *cfgValues) init() error {
  return c.parseCfgFile()
}

func (c *cfgValues) parseCfgFile() error {
  buf := [4096]byte{}
  for {
    n, err := c.buf.Read(buf[:])
    if err != nil {
      if err.Error() == "EOF" {
        break
      }
      return err
    }

    err = c.scan(buf[:n])
    if err != nil {
      return err
    }
  }

  return nil
}

func (cve *cfgValues) scan(buf []byte) error {
  _l := len(buf)
  var c byte
  // comm := ""
  // tempCfg := ""
  isCommProc := false
  for i := 0; i < _l; i++ {
    c = buf[i]
    switch c {
    case '\n':
      //_cfg += tempCfg
      cve.j += " "
      isCommProc = false
    case '\t':
      cve.j += " "
      continue
    case ';':
      //_cfg += tempCfg
      isCommProc = true
    default:
      if !unicode.IsPrint(rune(c)) || isCommProc {
        continue
      }

      cve.j += string(c)
    }
  }

  return nil
}
