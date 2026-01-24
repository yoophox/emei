package pki

//@key:[any space]  {type:<key type(ed25519/rsa...)>}.{cid:<id>}.{pri-len:<byts num>}.{pub-len:}\n
//<prikey>\n<pubkey>\n
//can't have any other letter before @key

// pem format
// public key:pem of PKIX, ASN.1 DER
// private key: pem of PKCS #8, ASN.1 DER

import (
  "bufio"
  "encoding/pem"
  "fmt"
  "io"
  "os"
  "regexp"
  "strconv"

  "github.com/fsnotify/fsnotify"
  "github.com/yoophox/emei/errs"
)

type local struct {
  cryptos map[uint64]*cryptosx
  crps    map[string][]*cryptosx // key = type: rsa
}

// newLocal ...
func newLocalBackend() *local {
  _, _ = fsnotify.NewWatcher()
  f, err := os.Open(_pkiLocal)
  if err != nil {
    panic(err)
  }
  defer f.Close()
  regS := regexp.MustCompile(`@key.*\{type:.*\}\.\{cid:\.*\}\.\{pri-len:.*\}\.\{pub-len:.*\}`)
  fkt := regexp.MustCompile(`\{type:.*\}`)
  fkid := regexp.MustCompile(`\{cid:.*\}`)
  fki := regexp.MustCompile(`\{pri-len:.*\}`)
  fku := regexp.MustCompile(`\{\pub-len:.*\}`)

  fsc := bufio.NewReader(f)
  var c *cryptosx
  var pril, publ int

  l := &local{cryptos: map[uint64]*cryptosx{}, crps: map[string][]*cryptosx{}}
  for line, err := fsc.ReadBytes('\n'); ; line, err = fsc.ReadBytes('\n') {
    if err != nil {
      if err == io.EOF {
        break
      }

      panic(err)
    }
    if len(line) == 0 {
      continue
    }

    if !regS.Match([]byte(line)) {
      panic("asfadadfas")
    }
    ct_ := fkt.Find([]byte(line))
    cid := fkid.Find([]byte(line))
    ci_ := fki.Find([]byte(line))
    cu_ := fku.Find([]byte(line))
    if string(ct_) == "" || string(ci_) == "" || string(cu_) == "" {
      panic("format error")
    }

    c = &cryptosx{}
    c.Typ = string(ct_[6 : len(ct_)-1])
    cid_, err := strconv.ParseUint(string(cid[5:len(cid)-1]), 10, 64)
    if err != nil {
      panic("parse cid")
    }
    c.CID = cid_
    pril, err = strconv.Atoi(string(ci_[9 : len(ci_)-1]))
    if err != nil {
      panic("parse pril")
    }
    publ, err = strconv.Atoi(string(cu_[9 : len(ci_)-1]))
    if err != nil {
      panic("parse publ")
    }
    ksc := make([]byte, pril+1+publ+1)
    n, err := fsc.Read(ksc)
    if n != (pril + 1 + publ + 1) {
      panic(n)
    }
    if err != nil {
      panic(err)
    }
    // var prib, pubb pem.Block
    prib, rest := pem.Decode(ksc[:pril])
    if rest != nil || len(rest) != 0 {
      panic(rest)
    }
    pubb, rest := pem.Decode(ksc[pril+1 : len(ksc)-1])
    if rest != nil || len(rest) != 0 {
      panic(rest)
    }
    c.Pri = prib
    c.Pub = pubb
    _, ok := l.crps[c.Typ]
    if !ok {
      l.crps[c.Typ] = []*cryptosx{}
    }
    l.crps[c.Typ] = append(l.crps[c.Typ], c)
    l.cryptos[cid_] = c

    c = nil
    pril = 0 //
    publ = 0
  }

  return l
}

func (l *local) getRandomCrypto(o ...Option) (uint64, error) {
  if len(l.cryptos) <= 0 {
    return 0, errs.Wrap(fmt.Errorf("have no keys in local data"), ERR_ID_PKI_LOCAL_NO_KEY)
  }

  return 0, nil
}
func (l *local) getPriKeyByID(id uint64) (*pem.Block, error)
func (l *local) getPubKeyByID(id uint64) (*pem.Block, error)
