package jwt

import (
  "crypto"
  "encoding/json"
  "fmt"
  "os"
  "time"

  "github.com/golang-jwt/jwt/v5"
  "github.com/yolksys/emei/cfg"
  "github.com/yolksys/emei/utils"
)

type Jwt interface {
  GetClaims(subject string, c any) error
  Err() error
}

// GenJwtStr ...
func GenJwtStr(subject string, c any, expire ...time.Duration) (string, error) {
  v, err := json.Marshal(c)
  if err != nil {
    return "", err
  }

  c_ := &claims{}
  c_.values = string(v)
  c_.Subject = subject
  tk_ := jwt.NewWithClaims(jwt.SigningMethodEdDSA, c_)
  return tk_.SignedString(_edPriKey)
}

// Parse ...
func Parse(jwtStr string) Jwt {
  token, err := jwt.ParseWithClaims(jwtStr, &claims{}, func(t *jwt.Token) (any, error) {
    return nil, nil
  })
  if token == nil {
    token = &jwt.Token{Header: map[string]any{}}
  } else if token.Header == nil {
    token.Header = map[string]any{}
  }
  if err != nil {
    token.Header["err"] = err
  }
  return (*jwtImpl)(token)
}

type claims struct {
  values string
  jwt.RegisteredClaims
}

type jwtImpl jwt.Token

func (j *jwtImpl) GetClaims(subject string, c any) error {
  t := (*jwt.Token)(j)
  s, err := t.Claims.GetSubject()
  if err != nil {
    return err
  }

  if s != subject {
    return fmt.Errorf("fail:GetClaims, reason:suject is not equal, recved:%s, req:%s", s, subject)
  }

  if t.Claims.(*claims).values == "" {
    return fmt.Errorf("fail:GetClaims, reason:have no values")
  }
  return json.Unmarshal([]byte(t.Claims.(*claims).values), c)
}

func (j *jwtImpl) Err() error {
  t := (*jwt.Token)(j)
  err, _ := t.Header["err"]
  if err != nil {
    return err.(error)
  }

  return nil
}

func init() {
  var path string

  err := cfg.GetCfgItem("edprikey", &path)
  utils.AssertErr(err)
  cnt, err := os.ReadFile(path)
  utils.AssertErr(err)
  _edPriKey, err = jwt.ParseEdPrivateKeyFromPEM(cnt)
  utils.AssertErr(err)

  err = cfg.GetCfgItem("edpubkey", &path)
  utils.AssertErr(err)
  cnt, err = os.ReadFile(path)
  utils.AssertErr(err)
  _edPriKey, err = jwt.ParseEdPublicKeyFromPEM(cnt)
  utils.AssertErr(err)
}

var (
  _edPriKey crypto.PrivateKey
  _edPunKey crypto.PublicKey
)
