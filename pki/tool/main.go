package main

import (
  "crypto/ed25519"
  "crypto/rand"
  "crypto/rsa"
  "crypto/x509"
  "encoding/pem"
  "flag"
  "fmt"
  "os"
  "time"
)

func main() {
}

func init() {
  fs_ := flag.NewFlagSet("pki-tool", flag.PanicOnError)
  en_ := fs_.Int("pki-tool-ed25519-num", 30, "ed25519 nums to be generated")
  ren_ := fs_.Int("pki-tool-rsa-num", 30, "rsa nums to be generated")
  fileP := fs_.String("pki-tool-file-path", "./pki.data", "out file path")
  fs_.Parse(os.Args[1:])
  _ed25519Num = *en_
  _rsaNum = *ren_

  f, err := os.OpenFile(*fileP, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
  if err != nil {
    panic(err)
  }

  for range _rsaNum {
    ct_, key, err := genRsaPem()
    if err != nil {
      panic(err)
    }

    write(f, "rsa", string(ct_), string(key))
  }

  for range _ed25519Num {
    pri, pub, err := genEd25519PEM()
    if err != nil {
      panic(err)
    }
    write(f, "ed25519", pri, pub)
  }
}

// write ...
func write(f *os.File, typ string, pri, pub string) {
  sepLine := "@key:   {type:%s}.{cid:%d}.{pri-len:%d}.{pub-len:%d}\n"
  sepLine = fmt.Sprintf(sepLine, typ, _startCID, len(pri), len(pub))
  _startCID++
  f.Write([]byte(sepLine))
  f.Write([]byte(pri))
  f.Write([]byte("\n"))
  f.Write([]byte(pub))
  f.Write([]byte("\n"))
}

// genEd25519PEM ...
func genEd25519PEM() (pri, pub string, err error) {
  pu_, p_, err := ed25519.GenerateKey(nil)
  if err != nil {
    return "", "", err
  }

  pd_, err := x509.MarshalPKCS8PrivateKey(p_)
  if err != nil {
    return "", "", err
  }
  pud, err := x509.MarshalPKIXPublicKey(pu_)
  if err != nil {
    return "", "", err
  }

  pri = string(pem.EncodeToMemory(&pem.Block{
    Type:  "ed25519 pri key",
    Bytes: pd_,
  }))

  pub = string(pem.EncodeToMemory(&pem.Block{
    Type:  "ed25519 pub key",
    Bytes: pud,
  }))

  return
}

// genRsa ...
func genRsaPem() (pri, pub string, err error) {
  privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
  if err != nil {
    return "", "", err
  }

  pd_, err := x509.MarshalPKCS8PrivateKey(privateKey)
  if err != nil {
    return "", "", err
  }
  pud, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
  if err != nil {
    return "", "", err
  }

  pri = string(pem.EncodeToMemory(&pem.Block{
    Type:  "rsa pri key",
    Bytes: pd_,
  }))

  pub = string(pem.EncodeToMemory(&pem.Block{
    Type:  "rsa pub key",
    Bytes: pud,
  }))

  return
}

var (
  _ed25519Num = 0
  _rsaNum     = 0
  _startCID   = time.Now().UnixNano()
)
